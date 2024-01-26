<script lang="ts" setup>
import { reactive, ref, nextTick, watch, onMounted } from "vue"
import {
  createArticleDataApi,
  deleteArticleDataApi,
  updateArticleDataApi,
  getArticleDataApi,
  getArticleDetailDataApi
} from "@/api/table/article"
import { CreateOrUpdateArticleRequestData, type GetArticleData } from "@/api/table/types/article"
import { type FormInstance, type FormRules, ElMessage, ElMessageBox } from "element-plus"
import { Search, Refresh, CirclePlus, Delete, RefreshRight } from "@element-plus/icons-vue"
import { usePagination } from "@/hooks/usePagination"
import { GetCollectionData } from "@/api/table/types/collection"
import { GetCategoryData } from "@/api/table/types/category"
import { GetTagData } from "@/api/table/types/tag"
import { getCollectionDataApi } from "@/api/table/collection"
import { getCategoryDataApi } from "@/api/table/category"
import { getTagDataApi } from "@/api/table/tag"

defineOptions({
  // 命名当前组件
  name: "Document"
})

const loading = ref<boolean>(false)
const { paginationData, handleCurrentChange, handleSizeChange } = usePagination()
const collectionData = ref<GetCollectionData[]>([])
const categoryData = ref<GetCategoryData[]>([])
const tagData = ref<GetTagData[]>([])

const getCollectionData = () => {
  loading.value = true
  getCollectionDataApi({})
    .then((res) => {
      collectionData.value = res.data
    })
    .catch(() => {
      tableData.value = []
    })
    .finally(() => {
      loading.value = false
    })
}

const getCategoryData = () => {
  loading.value = true
  getCategoryDataApi({})
    .then((res) => {
      categoryData.value = res.data
    })
    .catch(() => {
      tableData.value = []
    })
    .finally(() => {
      loading.value = false
    })
}

const getTagData = () => {
  loading.value = true
  getTagDataApi({})
    .then((res) => {
      tagData.value = res.data
    })
    .catch(() => {
      tableData.value = []
    })
    .finally(() => {
      loading.value = false
    })
}

const getArticleDetailData = (id: string) => {
  loading.value = true
  getArticleDetailDataApi(id)
    .then((res) => {
      formData.title = res.data.title
      formData.type = res.data.type
      formData.collectionId = res.data.collectionId == 0 ? undefined : res.data.collectionId
      formData.categoryId = res.data.categoryId == 0 ? undefined : res.data.categoryId
      formData.tags = res.data.tags
      formData.content = res.data.content
    })
    .finally(() => {
      loading.value = false
    })
}

//#region 增
const dialogVisible = ref<boolean>(false)
const dialogFullscreen = ref<boolean>(true)
const formRef = ref<FormInstance | null>(null)
const markdownRef = ref<FormInstance | null>(null)
const formData = ref<CreateOrUpdateArticleRequestData>({
  title: "",
  content: "",
  type: "",
  tags: [],
  collectionId: undefined,
  categoryId: undefined
}).value
const formRules: FormRules = reactive({
  title: [{ required: true, trigger: "blur", message: "请输入标题" }],
  type: [{ required: true, trigger: "blur", message: "请选择类别" }],
  content: [{ required: true, trigger: "blur", message: "请输入内容" }]
})
const handleCreateOrUpdate = () => {
  formRef.value?.validate((valid: boolean, fields) => {
    if (!valid) return console.error("表单校验不通过", fields)
    loading.value = true
    const api = currentUpdateId.value === undefined ? createArticleDataApi : updateArticleDataApi
    api({
      id: currentUpdateId.value,
      ...formData
    })
      .then(() => {
        ElMessage.success("操作成功")
        dialogVisible.value = false
        getArticleData()
      })
      .finally(() => {
        loading.value = false
      })
  })
}
const resetForm = () => {
  currentUpdateId.value = undefined
  formRef.value?.resetFields()
  markdownRef.value?.resetFields()
}
//#endregion

//#region 删
const handleDelete = (row: GetArticleData) => {
  ElMessageBox.confirm(`正在删除文章：${row.title}，确认删除？`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  }).then(() => {
    deleteArticleDataApi(row.ID).then(() => {
      ElMessage.success("删除成功")
      getArticleData()
    })
  })
}

const handleBatchDelete = () => {
  if (multipleSelection.value.length === 0) {
    // 如果 multipleSelection 是空数组，不执行删除操作
    ElMessage({
      message: "没有选中任何文章",
      type: "warning"
    })
    return
  }
  ElMessageBox.confirm(`正在批量删除选中文章，确认删除？`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  }).then(() => {
    // 创建一个数组来存储所有的删除请求
    const deletePromises = multipleSelection.value.map((item) => {
      // 返回每个删除请求的 Promise
      return deleteArticleDataApi(item.ID)
    })

    // 使用 Promise.all 等待所有请求完成
    Promise.all(deletePromises).then(() => {
      ElMessage.success("删除成功")
      // 全部删除成功后刷新列表
      getArticleData()
    })
  })
}
//#endregion

//#region 改
const currentUpdateId = ref<undefined | string>(undefined)
const handleUpdate = (row: GetArticleData) => {
  dialogVisible.value = true
  // 必须延迟赋值，防止 resetFields 方法将数据重置错误
  nextTick(() => {
    currentUpdateId.value = row.ID
    getArticleDetailData(row.ID)
  })
}
//#endregion

//#region 查
const tableData = ref<GetArticleData[]>([])
const searchFormRef = ref<FormInstance | null>(null)
const searchData = reactive({
  name: "",
  sort: "",
  type: ""
})
const typeOptions = reactive([
  { key: "doc", alias_name: "文档" },
  { key: "blog", alias_name: "博客" }
])
const sortOptions = reactive([
  { key: "-id", alias_name: "ID DESC" },
  { key: "+id", alias_name: "ID ASC" },
  { key: "-updated_at", alias_name: "UpdatedAt Desc" },
  { key: "+updated_at", alias_name: "UpdatedAt ASC" }
])
const getArticleData = () => {
  loading.value = true
  getArticleDataApi({
    keyword: searchData.name || undefined,
    sort: searchData.sort || undefined,
    type: searchData.type || undefined
  })
    .then((res) => {
      paginationData.total = res.data.total
      tableData.value = res.data.data
    })
    .catch(() => {
      tableData.value = []
    })
    .finally(() => {
      loading.value = false
    })
}
const handleSearch = () => {
  getArticleData()
}
const resetSearch = () => {
  searchFormRef.value?.resetFields()
  handleSearch()
}
const formatterDate = (row: GetArticleData) => {
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

const formatterType = (row: GetArticleData) => {
  return row.type == "doc" ? "文档" : "博客"
}
//#endregion

// 批量操作
const multipleSelection = ref<GetArticleData[]>([])
const handleSelectionChange = (val: GetArticleData[]) => {
  multipleSelection.value = val
}

/** 监听分页参数的变化 */
watch([() => paginationData.currentPage, () => paginationData.pageSize], getArticleData, { immediate: true })
/** 组件初始化加载数据 */
onMounted(() => {
  // 在组件挂载后执行的逻辑，例如加载数据
  getTagData()
  getCollectionData()
  getCategoryData()
})
</script>

<template>
  <div class="app-container">
    <el-card v-loading="loading" shadow="never" class="search-wrapper">
      <el-form ref="searchFormRef" :inline="true" :model="searchData">
        <el-form-item prop="name" label="标题">
          <el-input v-model="searchData.name" placeholder="请输入" />
        </el-form-item>
        <el-form-item prop="type" label="类别">
          <el-select v-model="searchData.type" placeholder="类别" style="width: 240px">
            <el-option v-for="item in typeOptions" :key="item.key" :label="item.alias_name" :value="item.key" />
          </el-select>
        </el-form-item>
        <el-form-item prop="sort" label="排序">
          <el-select v-model="searchData.sort" placeholder="排序字段" style="width: 240px">
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
          <el-button type="primary" :icon="CirclePlus" @click="dialogVisible = true">新增文章</el-button>
          <el-button type="danger" :icon="Delete" @click="handleBatchDelete">批量删除</el-button>
        </div>
        <div>
          <el-tooltip content="刷新当前页">
            <el-button type="primary" :icon="RefreshRight" circle @click="getArticleData" />
          </el-tooltip>
        </div>
      </div>
      <div class="table-wrapper">
        <el-table :data="tableData" @selection-change="handleSelectionChange">
          <el-table-column type="selection" width="50" align="center" />
          <el-table-column prop="title" label="标题" align="center" />
          <el-table-column prop="type" label="类别" align="center" :formatter="formatterType" />
          <el-table-column prop="categoryName" label="分类" align="center" />
          <el-table-column prop="collectionName" label="文集" align="center" />
          <el-table-column prop="UpdatedAt" label="修改时间" align="center" :formatter="formatterDate" />
          <el-table-column fixed="right" label="操作" width="150" align="center">
            <template #default="scope">
              <el-button type="primary" text bg size="small" @click="handleUpdate(scope.row)">修改</el-button>
              <el-button type="danger" text bg size="small" @click="handleDelete(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="pager-wrapper">
        <el-pagination
          background
          :layout="paginationData.layout"
          :page-sizes="paginationData.pageSizes"
          :total="paginationData.total"
          :page-size="paginationData.pageSize"
          :currentPage="paginationData.currentPage"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    <!-- 新增/修改 -->
    <el-dialog
      v-model="dialogVisible"
      :title="currentUpdateId === undefined ? '新增文章' : '修改文章'"
      @closed="resetForm"
      :fullscreen="dialogFullscreen"
    >
      <el-form
        :inline="true"
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        label-position="right"
        class="inline-form"
        style="width: 90%"
      >
        <el-form-item prop="title" label="标题">
          <el-input v-model="formData.title" placeholder="请输入" />
        </el-form-item>
        <el-form-item prop="type" label="类别">
          <el-select v-model="formData.type" placeholder="类别">
            <el-option v-for="item in typeOptions" :key="item.key" :label="item.alias_name" :value="item.key" />
          </el-select>
        </el-form-item>
        <el-form-item prop="categoryId" label="分类">
          <el-select v-model="formData.categoryId" placeholder="分类">
            <el-option v-for="item in categoryData" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item prop="collectionId" label="文集">
          <el-select v-model="formData.collectionId" placeholder="文集">
            <el-option v-for="item in collectionData" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item prop="tags" label="标签">
          <el-select multiple filterable allow-create default-first-option v-model="formData.tags" placeholder="标签">
            <el-option v-for="item in tagData" :key="item.ID" :label="item.name" :value="item.name" />
          </el-select>
        </el-form-item>
      </el-form>
      <el-form
        ref="markdownRef"
        :rules="formRules"
        label-width="100px"
        label-position="right"
        class="inline-form"
        style="width: 90%"
      >
        <el-form-item label="内容" prop="content">
          <mavon-editor v-model="formData.content" class="markdown-editor" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateOrUpdate" :loading="loading">确认</el-button>
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

.inline-form .el-input {
  --el-input-width: 480px;
}

.inline-form .el-select {
  --el-select-width: 480px;
}
.markdown-editor {
  width: 100%;
}
</style>
