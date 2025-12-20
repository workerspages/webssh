<template>
  <div class="system-login-container">
    <div class="login-box">
      <div class="logo">
        <i class="el-icon-lock" style="font-size: 48px; color: #409EFF;"></i>
      </div>
      <h2 class="title">WebSSH 管理系统</h2>
      <p class="subtitle">安全 · 高效 · 便捷</p>
      
      <el-form :model="form" ref="form" :rules="rules" class="login-form">
        <el-form-item prop="username">
          <el-input 
            v-model="form.username" 
            prefix-icon="el-icon-user" 
            placeholder="请输入管理员账号"
            size="large">
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input 
            v-model="form.password" 
            prefix-icon="el-icon-lock" 
            type="password" 
            placeholder="请输入密码" 
            show-password 
            size="large"
            @keyup.enter.native="handleLogin">
          </el-input>
        </el-form-item>
        <el-button type="primary" class="login-btn" :loading="loading" @click="handleLogin" size="large">立即登录</el-button>
      </el-form>
      
      <div class="footer-tip">默认账号: admin / 密码: admin123</div>
    </div>
  </div>
</template>

<script>
import request from '@/utils/request'

export default {
  name: 'SystemLogin',
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
              // 登录成功后跳转到首页（Layout）
              this.$router.push('/')
            }
          } catch (e) {
            console.error(e)
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
.system-login-container {
  height: 100vh;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f0f2f5 url('https://gw.alipayobjects.com/zos/rmsportal/TVYTbAXWheQpRcWDaDMu.svg');
  background-size: cover;
}

.login-box {
  background: rgba(255, 255, 255, 0.95);
  width: 420px;
  padding: 50px 40px;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.08);
  text-align: center;
  backdrop-filter: blur(10px);
}

.logo { margin-bottom: 10px; }
.title { margin: 10px 0; font-size: 26px; color: #303133; font-weight: 600; }
.subtitle { margin: 0 0 40px; color: #909399; font-size: 14px; letter-spacing: 2px; }

.login-form {
  text-align: left;
  .el-input__inner { height: 45px; }
}

.login-btn {
  width: 100%;
  font-size: 16px;
  padding: 12px 0;
  margin-top: 20px;
  letter-spacing: 4px;
}

.footer-tip {
  margin-top: 30px;
  font-size: 12px;
  color: #c0c4cc;
}
</style>
