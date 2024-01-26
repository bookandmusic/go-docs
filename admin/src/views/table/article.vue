<!-- eslint-disable vue/html-self-closing -->
<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.keyword" placeholder="标题" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.type" placeholder="请选择类别" clearable class="filter-item" style="width: 130px" @change="handleFilter">
        <el-option v-for="item in docTypeOptions" :key="item.key" :label="item.display_name" :value="item.key" />
      </el-select>
      <el-select v-model="listQuery.sort" style="width: 180px" placeholder="请选择排序字段" clearable class="filter-item" @change="handleFilter">
        <el-option v-for="item in sortOptions" :key="item.key" :label="item.label" :value="item.key" />
      </el-select>
      <el-button v-waves class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-search" @click="handleFilter">
        Search
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
        Add
      </el-button>
      <el-button v-waves :loading="downloadLoading" class="filter-item" type="primary" icon="el-icon-download" @click="handleDownload">
        Export
      </el-button>
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
      @sort-change="sortChange"
    >
      <el-table-column label="ID" prop="id" sortable="custom" align="center" width="80" :class-name="getSortClass('id')">
        <template slot-scope="{row}">
          <span>{{ row.ID }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Date" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.UpdatedAt | parseTime('{y}-{m}-{d} {h}:{i}') }}</span>
        </template>
      </el-table-column>
      <el-table-column label="标题" min-width="100px" prop="title" sortable="custom" :class-name="getSortClass('title')">
        <template slot-scope="{row}">
          <span class="link-type" @click="handleUpdate(row)">{{ row.title }}</span>
        </template>
      </el-table-column>
      <el-table-column label="标签数" width="110px" align="center">
        <template slot-scope="{row}">
          <span v-if="row.tagsNum" class="link-type" @click="handleDisplayTag(row.tags)">{{ row.tagsNum }}</span>
          <span v-else>0</span>
        </template>
      </el-table-column>
      <el-table-column label="类别" width="110px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.type }}</span>
        </template>
      </el-table-column>
      <el-table-column label="分类" width="110px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.categoryName }}</span>
        </template>
      </el-table-column>
      <el-table-column label="文集" width="110px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.collectionName }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Actions" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row,$index}">

          <el-button size="mini" type="danger" @click="handleDelete(row,$index)">
            Delete
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible" :fullscreen="dialogFormFullscreen">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="70px" style="width: 90%; margin-left:50px;">
        <el-form-item label="类别" prop="type">
          <el-select v-model="temp.type" placeholder="请选择" class="edit-form-element">
            <el-option v-for="item in docTypeOptions" :key="item.key" :label="item.display_name" :value="item.key"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="temp.title" class="edit-form-element" />
        </el-form-item>
        <el-form-item label="标签" prop="tags">
          <el-select v-model="temp.tags" multiple allow-create filterable clearable placeholder="请选择标签" class="edit-form-element">
            <el-option v-for="item in tags" :key="item.ID" :label="item.name" :value="item.name"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="temp.categoryId" clearable placeholder="请选择文章分类" class="edit-form-element">
            <el-option v-for="item in categories" :key="item.ID" :label="item.name" :value="item.ID"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="文集" prop="collection">
          <el-select v-model="temp.collectionId" clearable placeholder="请选择文集" class="edit-form-element">
            <el-option v-for="item in collections" :key="item.ID" :label="item.name" :value="item.ID"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <mavon-editor v-model="temp.content" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          Cancel
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          Confirm
        </el-button>
      </div>
    </el-dialog>
    <el-dialog :visible.sync="dialogTagsVisible" title="标签">
      <el-button v-for="tag, index in tagsData" :key="index" type="info"> {{ tag }} </el-button>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="dialogTagsVisible = false">Confirm</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { getArticleList, getArticleDetail, editArticle, deleteArticle } from '@/api/article'
import { getTagList } from '@/api/tag'
import { getCategoryList } from '@/api/category'
import { getCollectionList } from '@/api/collection'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import { mavonEditor } from 'mavon-editor'
import 'mavon-editor/dist/css/index.css'

export default {
  name: 'Article',
  components: { Pagination, mavonEditor },
  directives: { waves },
  filters: {
    parseTime
  },
  data() {
    return {
      tags: [],
      categories: [],
      collections: [],
      tableKey: 0,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        type: '',
        page: 1,
        limit: 20,
        keyword: undefined,
        sort: '-id'
      },
      importanceOptions: [1, 2, 3],
      docTypeOptions: [{ display_name: '博客', key: 'blog' }, { display_name: '文档', key: 'doc' }],
      sortOptions: [{ label: 'ID Ascending', key: '+id' }, { label: 'ID Descending', key: '-id' }, { label: 'Title Ascending', key: '+title' }, { label: 'Title Descending', key: '-title' }],
      temp: {
        id: 0,
        title: '',
        categoryId: 0,
        collectionId: 0,
        content: '',
        tags: [],
        type: 'blog'
      },
      dialogFormVisible: false,
      dialogFormFullscreen: true,
      dialogTagsVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑',
        create: '添加'
      },
      dialogPvVisible: false,
      tagsData: [],
      rules: {
        type: [{ required: true, message: '文章必须指定类别', trigger: 'blur' }],
        title: [{ required: true, message: '文章必须填写标题', trigger: 'blur' }],
        content: [{ required: true, message: '文章必须填写标题', trigger: 'blur' }]
      },
      downloadLoading: false
    }
  },
  created() {
    this.getList()
    this.getCategories()
    this.getTags()
    this.getCollections()
  },
  methods: {
    getList() {
      this.listLoading = true
      getArticleList(this.listQuery).then(response => {
        this.list = response.data.data
        this.total = response.data.total

        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 500)
      })
    },
    getCategories() {
      this.listLoading = true
      getCategoryList().then(response => {
        this.categories = response.data

        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 500)
      })
    },
    getTags() {
      this.listLoading = true
      getTagList().then(response => {
        this.tags = response.data

        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 500)
      })
    },
    getCollections() {
      this.listLoading = true
      getCollectionList().then(response => {
        this.collections = response.data

        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 500)
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetTemp() {
      this.temp = {
        id: 0,
        title: '',
        categoryId: 0,
        collectionId: 0,
        content: '',
        tagIds: [],
        type: 'blog'
      }
    },
    handleCreate() {
      this.resetTemp()
      this.temp.categoryId = this.temp.categoryId === 0 ? null : this.temp.categoryId
      this.temp.collectionId = this.temp.collectionId === 0 ? null : this.temp.collectionId
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          editArticle(this.temp).then(response => {
            this.list.unshift(response.data)
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: '添加文章成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleUpdate(row) {
      getArticleDetail(row.ID).then(response => {
        this.temp = response.data
        this.temp.categoryId = this.temp.categoryId === 0 ? null : this.temp.categoryId
        this.temp.collectionId = this.temp.collectionId === 0 ? null : this.temp.collectionId
      })
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          editArticle(tempData).then(response => {
            console.log(response)
            const index = this.list.findIndex(v => v.id === this.temp.id)
            this.list.splice(index, 1, response.data)
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: '更新文章成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleDelete(row, index) {
      deleteArticle(row.ID).then(() => {
        this.list.splice(index, 1)
      })
    },
    handleDownload() {
      this.downloadLoading = true
      import('@/vendor/Export2Excel').then(excel => {
        const tHeader = ['ID', '更新时间', '标题', '类别', '分类', '文集']
        const filterVal = ['ID', 'UpdatedAt', 'title', 'type', 'category_name', 'collection_name']
        const data = this.formatJson(filterVal)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: '文章列表'
        })
        this.downloadLoading = false
      })
    },
    formatJson(filterVal) {
      return this.list.map(v => filterVal.map(j => {
        if (j === 'UpdatedAt') {
          return parseTime(v[j])
        } else {
          return v[j]
        }
      }))
    },
    getSortClass: function(key) {
      const sort = this.listQuery.sort.slice(1)
      if (sort === key) {
        return sort === `+${key}` ? 'ascending' : 'descending'
      } else {
        return ''
      }
    },
    sortChange(data) {
      const { prop, order } = data
      if (order === 'ascending') {
        this.listQuery.sort = `+${prop}`
      } else {
        this.listQuery.sort = `-${prop}`
      }
      this.handleFilter()
    },
    handleDisplayTag(tags) {
      this.tagsData = tags
      this.dialogTagsVisible = true
    }
  }
}
</script>
<style>
.edit-form-element {
  width: 400px
}
</style>
