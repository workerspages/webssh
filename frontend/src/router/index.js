import Vue from 'vue'
import Router from 'vue-router'

// 引入组件
import Login from '@/views/Login.vue'          // 新版系统登录页
import Layout from '@/views/Layout.vue'        // 侧边栏布局
import SshTerminal from '@/views/SshTerminal.vue' // 旧版 SSH 连接功能
import TerminalPage from '@/views/TerminalPage.vue' // xterm 终端页
import CronList from '@/views/CronList.vue'    // 定时任务页
import NotifyConfig from '@/views/NotifyConfig.vue' // 通知配置页

Vue.use(Router)

const router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/login',
      component: Login
    },
    {
      path: '/',
      component: Layout,
      redirect: '/ssh', // 默认跳到 SSH 连接页
      children: [
        {
          path: 'ssh',
          component: SshTerminal,
          meta: { title: 'SSH终端', icon: 'el-icon-monitor' }
        },
        {
          path: 'cron',
          component: CronList,
          meta: { title: '定时任务', icon: 'el-icon-time' }
        },
        {
          path: 'notify',
          component: NotifyConfig,
          meta: { title: '通知配置', icon: 'el-icon-message-solid' }
        }
      ]
    },
    {
      path: '/terminal',
      component: TerminalPage
    }
  ]
})

// 路由守卫：检查 Token
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  // 如果去往 terminal 页面（新窗口打开的），可能不需要系统登录token，但需要 ssh 参数
  // 这里我们只保护系统管理页面
  if (to.path === '/login' || to.path === '/terminal') {
    next()
  } else {
    if (!token) {
      next('/login')
    } else {
      next()
    }
  }
})

export default router
