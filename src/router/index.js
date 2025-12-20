import Vue from 'vue'
import Router from 'vue-router'
import Login from '@/views/Login.vue' // 新的登录页
import Layout from '@/views/Layout.vue' // 侧边栏布局
import SshTerminal from '@/views/SshTerminal.vue' // 原 SSH 连接页面
import TerminalPage from '@/views/TerminalPage.vue' // 实际终端页
import CronList from '@/views/CronList.vue'
import NotifyConfig from '@/views/NotifyConfig.vue'

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
      redirect: '/ssh',
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

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
