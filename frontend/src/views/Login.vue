<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo-area">
        <i class="el-icon-lock logo-icon"></i>
        <h2 class="app-title">WebSSH</h2>
        <p class="app-desc">集中式SSH连接管理平台</p>
      </div>
      
      <el-form :model="form" ref="form" :rules="rules" class="login-form">
        <el-form-item prop="username">
          <el-input 
            v-model="form.username" 
            prefix-icon="el-icon-user" 
            placeholder="用户名"
            class="custom-input">
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input 
            v-model="form.password" 
            prefix-icon="el-icon-lock" 
            type="password" 
            placeholder="密码" 
            show-password 
            class="custom-input"
            @keyup.enter.native="handleLogin">
          </el-input>
        </el-form-item>
        <el-button type="primary" class="login-btn" :loading="loading" @click="handleLogin">登 录</el-button>
      </el-form>
      
      <div class="footer">
        <span>WebSSH © 2025 | GitHub</span>
      </div>
    </div>
    
    <!-- Decorative circles -->
    <div class="circle circle-1"></div>
    <div class="circle circle-2"></div>
    <div class="circle circle-3"></div>
  </div>
</template>

<script>
import request from '@/utils/request'

export default {
  name: 'Login',
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
.login-container {
  min-height: 100vh;
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;
}

.login-card {
  width: 380px;
  background: rgba(255, 255, 255, 0.9);
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
  z-index: 10;
  backdrop-filter: blur(10px);
  text-align: center;
}

.logo-area {
  margin-bottom: 30px;
  
  .logo-icon {
    font-size: 48px;
    background: linear-gradient(to right, #f6d365 0%, #fda085 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    margin-bottom: 10px;
  }

  .app-title {
    font-size: 28px;
    color: #333;
    margin: 10px 0 5px;
    font-weight: 700;
  }

  .app-desc {
    color: #888;
    font-size: 14px;
    margin: 0;
  }
}

.login-form {
  .custom-input {
    ::v-deep .el-input__inner {
      height: 45px;
      line-height: 45px;
      border-radius: 8px;
      background: #f5f7fa;
      border: 1px solid #e4e7ed;
      &:focus {
        background: #fff;
        border-color: #764ba2;
      }
    }
  }
}

.login-btn {
  width: 100%;
  height: 45px;
  font-size: 16px;
  border-radius: 8px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border: none;
  margin-top: 10px;
  transition: transform 0.2s;
  
  &:hover {
    opacity: 0.9;
    transform: translateY(-1px);
  }
  
  &:active {
    transform: translateY(0);
  }
}

.footer {
  margin-top: 30px;
  font-size: 12px;
  color: #999;
}

/* Decorative background circles */
.circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.circle-1 {
  width: 300px;
  height: 300px;
  top: -50px;
  left: -50px;
}

.circle-2 {
  width: 500px;
  height: 500px;
  bottom: -100px;
  right: -100px;
}

.circle-3 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  left: 20%;
}
</style>
