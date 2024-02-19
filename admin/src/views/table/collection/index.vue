<script lang="ts" setup>
import { reactive, ref, nextTick, onMounted } from "vue"
import {
  createCollectionDataApi,
  deleteCollectionDataApi,
  updateCollectionDataApi,
  getCollectionDataApi,
  bantchDeleteCollectionDataApi,
  getCollectionTocListDataApi,
  updateCollectionTocListDataApi
} from "@/api/table/collection"
import { type GetCollectionData, TocList } from "@/api/table/types/collection"
import { type FormInstance, type FormRules, ElMessage, ElMessageBox } from "element-plus"
import { Search, Refresh, CirclePlus, Delete, RefreshRight } from "@element-plus/icons-vue"

defineOptions({
  // 命名当前组件
  name: "Collection"
})

const loading = ref<boolean>(false)

//#region 增
const dialogVisible = ref<boolean>(false)
const formRef = ref<FormInstance | null>(null)
const formData = reactive({
  name: "",
  author: ""
})
const formRules: FormRules = reactive({
  name: [{ required: true, trigger: "blur", message: "请输入文集名" }],
  author: [{ required: true, trigger: "blur", message: "请输入作者" }]
})
const handleCreateOrUpdate = () => {
  formRef.value?.validate((valid: boolean, fields) => {
    if (!valid) return console.error("表单校验不通过", fields)
    loading.value = true
    const api = currentUpdateId.value === undefined ? createCollectionDataApi : updateCollectionDataApi
    api({
      id: currentUpdateId.value,
      ...formData
    })
      .then(() => {
        ElMessage.success("操作成功")
        dialogVisible.value = false
        getCollectionData()
      })
      .finally(() => {
        loading.value = false
      })
  })
}
const resetForm = () => {
  currentUpdateId.value = "0"
  formRef.value?.resetFields()
}
//#endregion

//#region 删
const handleDelete = (row: GetCollectionData) => {
  ElMessageBox.confirm(`正在删除文集：${row.name}，确认删除？`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  }).then(() => {
    deleteCollectionDataApi(row.ID).then(() => {
      ElMessage.success("删除成功")
      getCollectionData()
    })
  })
}

const handleBatchDelete = () => {
  if (multipleSelection.value.length === 0) {
    // 如果 multipleSelection 是空数组，不执行删除操作
    ElMessage({
      message: "没有选中任何文集",
      type: "warning"
    })
    return
  }
  ElMessageBox.confirm(`正在批量删除选中文集，确认删除？`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  }).then(() => {
    // 创建一个数组来存储所有的id
    const ids = multipleSelection.value.map((item) => {
      return item.ID
    })

    bantchDeleteCollectionDataApi(ids).then(() => {
      ElMessage.success("删除成功")
      // 全部删除成功后刷新列表
      getCollectionData()
    })
  })
}
//#endregion

//#region 改
const currentUpdateId = ref<string>("0")
const handleUpdate = (row: GetCollectionData) => {
  dialogVisible.value = true
  // 必须延迟赋值，防止 resetFields 方法将数据重置错误
  nextTick(() => {
    currentUpdateId.value = row.ID
    formData.name = row.name
    formData.author = row.author
  })
}
//#endregion

//#region 文章排序
const dialogTocListVisible = ref<boolean>(false)
const tocList = ref<TocList>([])
const tocProps = {
  label: "title"
}
const handleTocSorted = (row: GetCollectionData) => {
  dialogTocListVisible.value = true
  // 必须延迟赋值，防止 resetFields 方法将数据重置错误
  nextTick(() => {
    currentUpdateId.value = row.ID
    getCollectionTocListDataApi(row.ID).then((res) => {
      tocList.value = res.data
    })
  })
}

const handleDrop = () => {
  console.log("tree drop:", tocList)
  // 更新 TocList 的 order
  const updatedTocList = updateTocList(tocList.value, 0)
  console.log("tree drop:", updatedTocList)
}

const updateTocList = (tocList: TocList, parentId: number) => {
  // 遍历 tocList 数组
  tocList.forEach((item, index) => {
    // 更新 parentId 和 order 属性
    item.parent_id = parentId
    item.order = index + 1

    // 如果当前项有子目录，则递归调用 updateTocList 更新子目录
    if (item.children && item.children.length > 0) {
      updateTocList(item.children, item.id)
    }
  })

  return tocList
}

const handleUpdateTocSorted = () => {
  const updatedTocList = updateTocList(tocList.value, 0)
  updateCollectionTocListDataApi(currentUpdateId.value, updatedTocList).then(() => {
    ElMessage.success("排序成功")
    dialogTocListVisible.value = false
  })
}
//#endregion

//#region 查
const tableData = ref<GetCollectionData[]>([])
const searchFormRef = ref<FormInstance | null>(null)
const searchData = reactive({
  name: "",
  sort: ""
})
const sortOptions = reactive([
  { key: "-id", alias_name: "ID DESC" },
  { key: "+id", alias_name: "ID ASC" },
  { key: "-num", alias_name: "Num Desc" },
  { key: "+num", alias_name: "Num ASC" }
])
const getCollectionData = () => {
  loading.value = true
  getCollectionDataApi({
    keyword: searchData.name || undefined,
    sort: searchData.sort || undefined
  })
    .then((res) => {
      tableData.value = res.data
    })
    .catch(() => {
      tableData.value = []
    })
    .finally(() => {
      loading.value = false
    })
}
const handleSearch = () => {
  getCollectionData()
}
const resetSearch = () => {
  searchFormRef.value?.resetFields()
  handleSearch()
}

const formatterDate = (row: GetCollectionData) => {
  const originalDate = new Date(row.UpdatedAt)

  // 使用 Intl.DateTimeFormat 对象创建日期格式化选项
  const options: Intl.DateTimeFormatOptions = {
    year: "2-digit",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false, // 使用 24 小时制
    timeZone: "Asia/Shanghai" // 指定目标时区
  }

  // 使用 toLocaleString 方法将日期转换为字符串
  const formattedDate = new Intl.DateTimeFormat("zh-CN", options).format(originalDate)

  return formattedDate
}
//#endregion

// 批量操作
const multipleSelection = ref<GetCollectionData[]>([])
const handleSelectionChange = (val: GetCollectionData[]) => {
  multipleSelection.value = val
}

/** 组件初始化加载数据 */
onMounted(() => {
  // 在组件挂载后执行的逻辑，例如加载数据
  getCollectionData()
})
</script>

<template>
  <div class="app-container">
    <el-card v-loading="loading" shadow="never" class="search-wrapper">
      <el-form ref="searchFormRef" :inline="true" :model="searchData">
        <el-form-item prop="name" label="文集名">
          <el-input v-model="searchData.name" placeholder="请输入" />
        </el-form-item>
        <el-form-item prop="sort" label="排序">
          <el-select v-model="searchData.sort" placeholder="Select" style="width: 240px">
            <el-option v-for="item in sortOptions" :key="item.key" :label="item.alias_name" :value="item.key" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    <el-card v-loading="loading" shadow="never">
      <div class="toolbar-wrapper">
        <div>
          <el-button type="primary" :icon="CirclePlus" @click="dialogVisible = true">新增文集</el-button>
          <el-button type="danger" :icon="Delete" @click="handleBatchDelete">批量删除</el-button>
        </div>
        <div>
          <el-tooltip content="刷新当前页">
            <el-button type="primary" :icon="RefreshRight" circle @click="getCollectionData" />
          </el-tooltip>
        </div>
      </div>
      <div class="table-wrapper">
        <el-table :data="tableData" @selection-change="handleSelectionChange">
          <el-table-column type="selection" width="50" align="center" />
          <el-table-column prop="name" label="文集名" align="center" />
          <el-table-column prop="author" label="作者名" align="center" />
          <el-table-column prop="num" label="文章数" align="center" />
          <el-table-column prop="UpdatedAt" label="修改时间" align="center" :formatter="formatterDate" />
          <el-table-column fixed="right" label="操作" width="300" align="center">
            <template #default="scope">
              <el-button type="primary" text bg size="small" @click="handleUpdate(scope.row)">修改</el-button>
              <el-button type="primary" text bg size="small" @click="handleTocSorted(scope.row)">排序</el-button>
              <el-button type="danger" text bg size="small" @click="handleDelete(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    <!-- 新增/修改 -->
    <el-dialog
      v-model="dialogVisible"
      :title="currentUpdateId === undefined ? '新增文集' : '修改文集'"
      @closed="resetForm"
      width="30%"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="left">
        <el-form-item prop="name" label="文集名">
          <el-input v-model="formData.name" placeholder="请输入" />
        </el-form-item>
        <el-form-item prop="author" label="作者">
          <el-input v-model="formData.author" placeholder="请输入" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateOrUpdate" :loading="loading">确认</el-button>
      </template>
    </el-dialog>
    <!-- 文章排序 -->
    <el-dialog v-model="dialogTocListVisible" title="文章列表" width="30%">
      <el-tree :data="tocList" draggable default-expand-all node-key="id" :props="tocProps" @node-drop="handleDrop" />
      <template #footer>
        <el-button @click="dialogTocListVisible = false">取消</el-button>
        <el-button type="primary" @click="handleUpdateTocSorted" :loading="loading">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style lang="scss" scoped>
.search-wrapper {
  margin-bottom: 20px;
  :deep(.el-card__body) {
    padding-bottom: 2px;
  }
}

.toolbar-wrapper {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.table-wrapper {
  margin-bottom: 20px;
}

.pager-wrapper {
  display: flex;
  justify-content: flex-end;
}
</style>
