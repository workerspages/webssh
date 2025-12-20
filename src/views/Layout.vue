<template>
  <el-container style="height: 100vh;">
    <el-header class="app-header">
      <div class="brand">
        <i class="el-icon-lock" style="color: #ff9800;"></i> WebSSH
      </div>
      <div class="user-info">
        <el-dropdown @command="handleCommand">
          <span class="el-dropdown-link" style="color: #fff; cursor: pointer;">
            管理员 <i class="el-icon-arrow-down el-icon--right"></i>
          </span>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item command="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </el-header>
    
    <el-container style="overflow: hidden;">
      <el-aside width="200px" style="background-color: #fff; border-right: 1px solid #eee;">
        <el-menu :default-active="$route.path" router style="border:none;">
          <el-menu-item index="/ssh">
            <i class="el-icon-monitor"></i>
            <span slot="title">SSH连接</span>
          </el-menu-item>
          <el-menu-item index="/cron">
            <i class="el-icon-time"></i>
            <span slot="title">定时任务</span>
          </el-menu-item>
          <el-menu-item index="/notify">
            <i class="el-icon-message-solid"></i>
            <span slot="title">通知配置</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-main style="background: #f5f7fa; padding: 20px;">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
export default {
  methods: {
    handleCommand(cmd) {
      if (cmd === 'logout') {
        localStorage.removeItem('token')
        this.$router.push('/login')
      }
    }
  }
}
</script>

<style scoped>
.app-header {
  background-color: #333;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
}
.brand { font-size: 20px; font-weight: bold; display: flex; align-items: center; gap: 10px; }
</style>
