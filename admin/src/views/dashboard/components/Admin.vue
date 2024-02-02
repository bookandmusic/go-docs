<template>
  <div class="app-container center">
    <el-row :gutter="12">
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <el-button size="large" color="#7FC34D" bg icon="Document" circle />
              <el-button class="tip" text type="success">文章统计</el-button>
            </div>
          </template>
          <div>
            <div class="bottom">
              <span>{{ articleInfo.article_count }}</span>
            </div>
            <div class="tip">
              <span>当前文章总记录数</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <el-button size="large" color="#409EFF" bg icon="CollectionTag" circle />
              <el-button class="tip" text type="primary">标签统计</el-button>
            </div>
          </template>
          <div>
            <div class="bottom">
              <span>{{ articleInfo.tag_count }}</span>
            </div>
            <div class="tip">
              <span>当前标签总记录数</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <el-button size="large" color="#E6A23C" bg icon="Collection" circle />
              <el-button class="tip" text type="warning">文集统计</el-button>
            </div>
          </template>
          <div>
            <div class="bottom">
              <span>{{ articleInfo.collection_count }}</span>
            </div>
            <div class="tip">
              <span>当前文集总记录数</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <el-button size="large" color="#F56E6C" bg icon="FolderOpened" circle />
              <el-button class="tip" text type="danger">分类统计</el-button>
            </div>
          </template>
          <div>
            <div class="bottom">
              <span>{{ articleInfo.category_count }}</span>
            </div>
            <div class="tip">
              <span>当前分类总记录数</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue"
import { type ArticleInfoData } from "@/api/dashboard/types/index"
import { getArticleInfoApi } from "@/api/dashboard/index"
//#region 查
const articleInfo = ref<ArticleInfoData>({
  article_count: 0,
  tag_count: 0,
  category_count: 0,
  collection_count: 0
})
const loading = ref<boolean>(false)

const getArticleInfo = () => {
  loading.value = true
  getArticleInfoApi()
    .then((res) => {
      articleInfo.value = res.data
    })
    .catch(() => {
      articleInfo.value = {}
    })
    .finally(() => {
      loading.value = false
    })
}
/** 组件初始化加载数据 */
onMounted(() => {
  // 在组件挂载后执行的逻辑，例如加载数据
  getArticleInfo()
})
</script>

<style lang="scss" scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.bottom {
  margin: 13px;
  font-size: 30px;
}
.tip {
  margin: 13px;
  font-size: 14px;
}
</style>
