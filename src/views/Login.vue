<template>
  <div class="login-page">
    <div class="login-box">
      <div class="logo">
        <i class="el-icon-lock" style="font-size: 40px; color: #ff9800;"></i>
      </div>
      <h2 class="title">WebSSH</h2>
      <p class="subtitle">安全的SSH终端管理平台</p>
      
      <el-form :model="form" ref="form" :rules="rules" class="login-form">
        <el-form-item prop="username">
          <el-input v-model="form.username" prefix-icon="el-icon-user" placeholder="用户名"></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" prefix-icon="el-icon-lock" type="password" placeholder="密码" show-password @keyup.enter.native="handleLogin"></el-input>
        </el-form-item>
        <el-button type="primary" class="login-btn" :loading="loading" @click="handleLogin">登 录</el-button>
      </el-form>
      
      <div class="footer">WebSSH © 2025 | GitHub</div>
    </div>
  </div>
</template>

<script>
import request from '@/utils/request'

export default {
  data() {
    return {
      form: { username: '', password: '' },
      rules: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      },
      loading: false
    }
  },
  methods: {
    handleLogin() {
      this.$refs.form.validate(async valid => {
        if (valid) {
          this.loading = true
          try {
            const res = await request.post('/login', this.form)
            if (res.code === 200) {
              localStorage.setItem('token', res.token)
              this.$message.success('登录成功')
              this.$router.push('/')
            }
          } catch (e) {
            // error handled in interceptor
          } finally {
            this.loading = false
          }
        }
      })
    }
  }
}
</script>

<style scoped lang="scss">
.login-page {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-box {
  background: #fff;
  width: 400px;
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 15px 30px rgba(0,0,0,0.1);
  text-align: center;
}
.logo { margin-bottom: 10px; }
.title { margin: 0; font-size: 24px; color: #333; }
.subtitle { margin: 10px 0 30px; color: #888; font-size: 14px; }
.login-btn { width: 100%; background: linear-gradient(to right, #667eea, #764ba2); border: none; }
.footer { margin-top: 30px; font-size: 12px; color: #999; }
</style>
