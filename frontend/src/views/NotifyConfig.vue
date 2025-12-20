<template>
  <div class="notify-config-container">
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>通知配置</span>
      </div>
      
      <el-tabs v-model="activeTab" type="card">
        <el-tab-pane label="邮件通知" name="email">
          <el-form :model="form" label-width="120px" class="config-form">
            <el-form-item label="启用">
              <el-switch v-model="form.enable_email"></el-switch>
            </el-form-item>
            <el-form-item label="SMTP服务器">
              <el-input v-model="form.email_host" placeholder="例如: smtp.gmail.com"></el-input>
            </el-form-item>
            <el-form-item label="SMTP端口">
              <el-input v-model.number="form.email_port" placeholder="例如: 587"></el-input>
            </el-form-item>
            <el-form-item label="邮箱账号">
              <el-input v-model="form.email_user" placeholder="发送者邮箱地址"></el-input>
            </el-form-item>
            <el-form-item label="邮箱密码">
              <el-input v-model="form.email_pass" type="password" show-password placeholder="邮箱密码或应用授权码"></el-input>
            </el-form-item>
            <el-form-item label="接收邮箱">
              <el-input v-model="form.email_to" placeholder="接收通知的邮箱"></el-input>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="Telegram通知" name="telegram">
          <el-form :model="form" label-width="120px" class="config-form">
            <el-form-item label="启用">
              <el-switch v-model="form.enable_tg"></el-switch>
            </el-form-item>
            <el-form-item label="Bot Token">
              <el-input v-model="form.tg_bot_token" placeholder="BotFather 提供的 Token"></el-input>
            </el-form-item>
            <el-form-item label="Chat ID">
              <el-input v-model="form.tg_chat_id" placeholder="接收消息的 Chat ID"></el-input>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="Bark通知" name="bark">
          <el-form :model="form" label-width="120px" class="config-form">
             <el-form-item label="启用">
              <el-switch v-model="form.enable_bark"></el-switch>
            </el-form-item>
            <el-form-item label="Bark URL">
              <el-input v-model="form.bark_url" placeholder="https://api.day.app/Key/ 或仅 Key"></el-input>
              <div style="font-size: 12px; color: #909399; line-height: 1.5">
                支持完整 API URL 或仅填写 Key。<br>
                如果是自建服务请填写完整 URL (结尾不要带 /)。
              </div>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div class="action-footer">
        <el-button type="primary" @click="saveConfig" :loading="saving">保存配置</el-button>
        <el-button type="success" @click="testNotify" :loading="testing">测试发送</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import request from '@/utils/request'

export default {
  name: 'NotifyConfig',
  data() {
    return {
      activeTab: 'email',
      form: {
        enable_email: false,
        email_host: '',
        email_port: 465,
        email_user: '',
        email_pass: '',
        email_to: '',
        enable_tg: false,
        tg_bot_token: '',
        tg_chat_id: '',
        enable_bark: false,
        bark_url: ''
      },
      saving: false,
      testing: false
    }
  },
  created() {
    this.fetchConfig()
  },
  methods: {
    async fetchConfig() {
      try {
        const res = await request.get('/notify/config')
        if (res.code === 200 && res.data.ID) {
          // Merge API data into form
          this.form = { ...this.form, ...res.data }
        }
      } catch (e) {
        console.error(e)
      }
    },
    async saveConfig() {
      this.saving = true
      try {
        const res = await request.post('/notify/save', this.form)
        if (res.code === 200) {
          this.$message.success('保存成功')
        }
      } catch (e) {
        console.error(e)
      } finally {
        this.saving = false
      }
    },
    async testNotify() {
      this.testing = true
      try {
        const res = await request.post('/notify/test')
        if (res.code === 200) {
          this.$message.success('测试消息已发送，请检查接收情况')
        }
      } catch (e) {
        console.error(e)
      } finally {
        this.testing = false
      }
    }
  }
}
</script>

<style scoped lang="scss">
.notify-config-container {
  padding: 20px;
  .box-card {
    max-width: 800px;
    margin: 0 auto;
  }
  .config-form {
    margin-top: 20px;
    max-width: 600px;
  }
  .action-footer {
    margin-top: 30px;
    padding-left: 120px;
  }
}
</style>
