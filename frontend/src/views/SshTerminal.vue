<template>
  <div class="ssh-connect-container">
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>快速连接 SSH</span>
        <el-button style="float: right; padding: 3px 0" type="text" @click="onReset">重置表单</el-button>
      </div>
      
      <el-form :model="sshInfo" label-position="top" class="form-grid">
        <el-row :gutter="20">
           <el-col :span="12">
             <el-form-item label="主机地址 (Hostname)">
               <el-input ref="hostnameInput" v-model="sshInfo.hostname" placeholder="请输入主机地址" prefix-icon="el-icon-monitor"/>
             </el-form-item>
           </el-col>
           <el-col :span="12">
             <el-form-item label="端口 (Port)">
               <el-input v-model.number="sshInfo.port" placeholder="22" prefix-icon="el-icon-connection"/>
             </el-form-item>
           </el-col>
         </el-row>
         
         <el-row :gutter="20">
           <el-col :span="12">
             <el-form-item label="用户名 (Username)">
               <el-input ref="usernameInput" v-model="sshInfo.username" placeholder="root" prefix-icon="el-icon-user"/>
             </el-form-item>
           </el-col>
           <el-col :span="12">
             <el-form-item label="密码 (Password)">
               <el-input ref="passwordInput" v-model="sshInfo.password" type="password" placeholder="请输入密码" show-password prefix-icon="el-icon-lock"/>
             </el-form-item>
           </el-col>
         </el-row>
         
         <el-row :gutter="20">
           <el-col :span="12">
             <el-form-item label="私钥 (Private Key)">
               <el-upload
                 class="upload-key"
                 action=""
                 :show-file-list="false"
                 :before-upload="handlePrivateKeyUpload"
                 accept=".pem,.ppk,.key,.rsa,.id_rsa,.id_dsa,.txt">
                 <el-button size="small" icon="el-icon-upload2">上传密钥文件</el-button>
                 <span v-if="privateKeyFileName" class="file-name">{{ privateKeyFileName }}</span>
               </el-upload>
             </el-form-item>
           </el-col>
           <el-col :span="12">
             <el-form-item label="密钥口令 (Passphrase)">
               <el-input v-model="sshInfo.passphrase" placeholder="如有请填写" type="password" prefix-icon="el-icon-key"/>
             </el-form-item>
           </el-col>
         </el-row>
         
        <el-row>
          <el-col :span="24">
            <el-form-item label="初始命令 (Initial command)">
              <el-input v-model="sshInfo.command" placeholder="登录后自动执行的命令，如: htop" prefix-icon="el-icon-cpu"/>
            </el-form-item>
          </el-col>
        </el-row>

        <div class="action-bar">
          <el-button type="primary" icon="el-icon-link" @click="onGenerateLink">生成快捷链接</el-button>
          <el-button type="success" @click="onConnect" class="connect-btn">
            <i class="el-icon-video-play"></i> 连接终端
          </el-button>
        </div>

        <el-row v-if="generatedLink" style="margin-top: 18px;">
          <el-col :span="24">
            <el-input v-model="generatedLink" readonly>
              <template slot="append">
                <el-button icon="el-icon-document-copy" @click="copyLink">复制</el-button>
              </template>
            </el-input>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'SshTerminal',
  data () {
    return {
      sshInfo: {
        hostname: '',
        port: 22,
        username: '',
        password: '',
        privateKey: '',
        passphrase: '',
        command: ''
      },
      privateKeyFileName: '',
      generatedLink: ''
    }
  },
  created() {
    const savedInfo = localStorage.getItem('connectionInfo')
    if (savedInfo) {
      try {
        const info = JSON.parse(savedInfo)
        this.sshInfo = { ...this.sshInfo, ...info }
        if (info.privateKey) this.privateKeyFileName = '已保存密钥'
      } catch(e) {/*ignore*/}
    }
  },
  methods: {
    onConnect () {
      if (!this.sshInfo.hostname || !this.sshInfo.username) {
        this.$message.error('主机地址和用户名不能为空')
        return
      }
      
      // 清理数据
      if (this.sshInfo.privateKey) {
        this.sshInfo.password = ''
      }

      // 保存配置
      localStorage.setItem('connectionInfo', JSON.stringify(this.sshInfo))
      
      // 构造跳转参数
      const query = {
        hostname: encodeURIComponent(this.sshInfo.hostname),
        port: this.sshInfo.port,
        username: encodeURIComponent(this.sshInfo.username),
        command: encodeURIComponent(this.sshInfo.command)
      }
      
      if (this.sshInfo.privateKey) {
        sessionStorage.setItem('sshInfo', JSON.stringify(this.sshInfo))
        query.useKey = 1
      } else {
        query.password = btoa(this.sshInfo.password)
      }

      // 打开新窗口
      const routeUrl = this.$router.resolve({ path: '/terminal', query })
      window.open(routeUrl.href, '_blank')
    },
    handlePrivateKeyUpload(file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        this.sshInfo.privateKey = e.target.result
        this.privateKeyFileName = file.name
        this.sshInfo.password = '' // 互斥
      }
      reader.readAsText(file)
      return false
    },
    onReset() {
      this.sshInfo = { port: 22, hostname: '', username: '', password: '', command: '', privateKey: '' }
      this.privateKeyFileName = ''
      this.generatedLink = ''
    },
    onGenerateLink() {
        // 生成链接逻辑同前，略...
        this.$message.info('生成链接功能暂略')
    },
    copyLink() {
       // 复制逻辑
    }
  }
}
</script>

<style scoped>
.ssh-connect-container {
  max-width: 800px;
  margin: 0 auto;
}
.file-name {
  margin-left: 10px;
  font-size: 12px;
  color: #67c23a;
}
.action-bar {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
  gap: 15px;
}
.connect-btn {
  padding-left: 30px;
  padding-right: 30px;
  font-weight: bold;
}
</style>
