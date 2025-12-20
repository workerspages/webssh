<template>
  <div class="app-wrapper">
    <el-container style="height: 100vh">
      <el-aside width="220px" class="sidebar-container">
        <div class="logo-wrapper">
          <i class="el-icon-lock logo-icon"></i>
          <span v-if="!isCollapse">WebSSH</span>
        </div>
        <el-menu
          :default-active="$route.path"
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
          :collapse="isCollapse"
          router
          class="el-menu-vertical">
          
          <el-menu-item index="/ssh">
            <i class="el-icon-monitor"></i>
            <span slot="title">SSH 终端</span>
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
        
        <div class="bottom-actions">
           <el-button type="text" @click="logout" style="color: #bfcbd9">
             <i class="el-icon-switch-button"></i> 退出登录
           </el-button>
        </div>
      </el-aside>
      
      <el-container>
        <el-header class="navbar">
          <div class="header-left">
             <i :class="isCollapse ? 'el-icon-s-unfold' : 'el-icon-s-fold'" @click="toggleCollapse" class="hamburger"></i>
             <el-breadcrumb separator="/" class="breadcrumb">
                <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
                <el-breadcrumb-item>{{ $route.meta.title }}</el-breadcrumb-item>
             </el-breadcrumb>
          </div>
        </el-header>
        
        <el-main class="app-main">
          <transition name="fade-transform" mode="out-in">
            <router-view />
          </transition>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
export default {
  name: 'Layout',
  data() {
    return {
      isCollapse: false
    }
  },
  methods: {
    toggleCollapse() {
      this.isCollapse = !this.isCollapse
    },
    logout() {
      localStorage.removeItem('token')
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped lang="scss">
.app-wrapper {
  position: relative;
  height: 100%;
  width: 100%;
}

.sidebar-container {
  background-color: #304156;
  transition: width 0.28s;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  
  .logo-wrapper {
    height: 60px;
    line-height: 60px;
    text-align: center;
    color: #fff;
    font-weight: 600;
    font-size: 20px;
    background: #2b3a4d;
    
    .logo-icon {
      margin-right: 5px;
      color: #409EFF;
    }
  }
  
  .el-menu-vertical {
    border-right: none;
    flex: 1;
  }
  
  .bottom-actions {
    padding: 20px;
    text-align: center;
    background: #263445;
  }
}

.navbar {
  height: 60px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
  display: flex;
  align-items: center;
  padding: 0 20px;
  
  .hamburger {
    font-size: 20px;
    cursor: pointer;
    margin-right: 15px;
  }
  
  .header-left {
    display: flex;
    align-items: center;
  }
}

.app-main {
  min-height: calc(100vh - 60px);
  width: 100%;
  position: relative;
  overflow: hidden;
  background-color: #f0f2f5;
  padding: 20px;
}

/* Transitions */
.fade-transform-leave-active,
.fade-transform-enter-active {
  transition: all .5s;
}

.fade-transform-enter {
  opacity: 0;
  transform: translateX(-30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
