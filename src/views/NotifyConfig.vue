<template>
  <div class="notify-container">
    <div class="header">
      <h2>通知配置</h2>
    </div>
    
    <el-tabs type="border-card">
      <el-tab-pane label="邮件通知">
        <el-form :model="form" label-width="120px" style="max-width: 600px; margin-top: 20px;">
          <el-form-item label="启用">
            <el-switch v-model="form.EnableEmail"></el-switch>
          </el-form-item>
          <el-form-item label="SMTP主机">
            <el-input v-model="form.EmailHost" placeholder="smtp.example.com"></el-input>
          </el-form-item>
          <el-form-item label="SMTP端口">
            <el-input-number v-model="form.EmailPort" :min="1" :max="65535"></el-input-number>
          </el-form-item>
          <el-form-item label="用户名">
            <el-input v-model="form.EmailUser" placeholder="your@email.com"></el-input>
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="form.EmailPass" type="password" placeholder="SMTP密码或授权码" show-password></el-input>
          </el-form-item>
          <el-form-item label="发件人">
            <el-input v-model="form.EmailUser" placeholder="默认与用户名相同" disabled></el-input>
          </el-form-item>
          <el-form-item label="收件人">
            <el-input v-model="form.EmailTo" placeholder="接收通知的邮箱地址"></el-input>
          </el-form-item>
          <el-form-item label="使用TLS">
             <el-checkbox checked disabled>启用TLS加密</el-checkbox>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      
      <el-tab-pane label="Telegram通知">
         <el-form :model="form" label-width="120px" style="max-width: 600px; margin-top: 20px;">
          <el-form-item label="启用">
            <el-switch v-model="form.EnableTg"></el-switch>
          </el-form-item>
          <el-form-item label="Bot Token">
            <el-input v-model="form.TgBotToken" placeholder="例如: 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"></el-input>
            <div class="tip">通过 @BotFather 创建Bot获取Token</div>
          </el-form-item>
          <el-form-item label="Chat ID">
            <el-input v-model="form.TgChatID" placeholder="例如: 123456789"></el-input>
            <div class="tip">发送消息给Bot后，通过 getUpdates API 获取</div>
          </el-form-item>
         </el-form>
      </el-tab-pane>
    </el-tabs>
    
    <div class="actions">
      <el-button type="primary" @click="save">保存配置</el-button>
      <el-button @click="test">测试发送</el-button>
    </div>
  </div>
</template>

<script>
import request from '@/utils/request'

export default {
  data() {
    return {
      form: {
        EnableEmail: false,
        EmailHost: '',
        EmailPort: 465,
        EmailUser: '',
        EmailPass: '',
        EmailTo: '',
        EnableTg: false,
        TgBotToken: '',
        TgChatID: ''
      }
    }
  },
  created() {
    this.load()
  },
  methods: {
    async load() {
      const res = await request.get('/notify/config')
      if(res.data.ID) this.form = res.data
    },
    async save() {
      await request.post('/notify/save', this.form)
      this.$message.success('保存成功')
    },
    async test() {
        await request.post('/notify/test')
        this.$message.success('测试消息已发送，请检查接收端')
    }
  }
}
</script>

<style scoped>
.notify-container { padding: 20px; }
.header { margin-bottom: 20px; }
.actions { margin-top: 20px; }
.tip { font-size: 12px; color: #999; }
</style>
