<template>
  <div class="cron-list-container">
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>定时任务管理</span>
        <el-button style="float: right; padding: 3px 0" type="text" icon="el-icon-plus" @click="handleAdd">新建任务</el-button>
      </div>
      
      <el-table :data="jobs" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="Name" label="任务名称" width="180"></el-table-column>
        <el-table-column prop="CronExpr" label="Cron表达式" width="150">
          <template slot-scope="scope">
            <el-tag size="small">{{ scope.row.CronExpr }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="Status" label="状态" width="100">
          <template slot-scope="scope">
            <el-tag :type="scope.row.Status === 1 ? 'success' : 'info'">
              {{ scope.row.Status === 1 ? '运行中' : '已停止' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="上次执行" width="180">
          <template slot-scope="scope">
             {{ formatTime(scope.row.LastRunTime) }}
          </template>
        </el-table-column>
         <el-table-column label="上次结果">
          <template slot-scope="scope">
             <el-popover trigger="hover" placement="top" width="400">
               <pre class="log-preview">{{ scope.row.ErrorLog || '无日志' }}</pre>
               <div slot="reference" class="name-wrapper">
                 <el-tag size="medium" :type="getResultType(scope.row.LastResult)">{{ scope.row.LastResult || '无' }}</el-tag>
               </div>
             </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250">
          <template slot-scope="scope">
            <el-button size="mini" type="primary" @click="handleRun(scope.row)">执行</el-button>
            <el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button size="mini" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Add/Edit Dialog -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="60%">
      <el-form :model="form" ref="form" label-width="120px" :rules="rules">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="基本信息" name="basic">
            <el-form-item label="任务名称" prop="Name" style="margin-top: 20px">
              <el-input v-model="form.Name" placeholder="任务名称"></el-input>
            </el-form-item>
            <el-form-item label="Cron表达式" prop="CronExpr">
              <el-input v-model="form.CronExpr" placeholder="例如: 20 07 * * * (每日早上7点20分)"></el-input>
              <div class="tip">
                格式: <span style="font-weight: bold; color: #409EFF">分 时 日 月 周</span> (5位标准格式)<br>
                参考: <b>20 07 * * *</b> &nbsp;&nbsp;(每日早上7点20分)<br>
                参考: <b>*/30 9-23 * * *</b> (早9点至23点, 每30分)<br>
                参考: <b>*/5 * * * *</b> &nbsp;&nbsp;&nbsp;(每5分钟)
              </div>
            </el-form-item>
            <el-form-item label="随机延迟">
              <el-input-number v-model="form.RandomDelay" :min="0" :max="1440" controls-position="right"></el-input-number>
              <span style="margin-left: 10px">分钟</span>
              <div class="tip">任务触发后，将在 0 ~ N 分钟内随机执行 (防检测)</div>
            </el-form-item>
            <el-form-item label="状态">
              <el-switch v-model="form.Status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="停用"></el-switch>
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane label="SSH配置" name="ssh">
            <el-form-item label="主机地址" required style="margin-top: 20px">
              <el-input v-model="sshForm.hostname" placeholder="IP地址或域名"></el-input>
            </el-form-item>
            <el-form-item label="端口">
              <el-input-number v-model="sshForm.port" :min="1" :max="65535"></el-input-number>
            </el-form-item>
            <el-form-item label="用户名" required>
              <el-input v-model="sshForm.username" placeholder="root"></el-input>
            </el-form-item>
            <el-form-item label="认证方式">
              <el-radio-group v-model="sshForm.logintype">
                <el-radio :label="0">密码</el-radio>
                <el-radio :label="1">密钥</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="密码" v-if="sshForm.logintype === 0">
              <el-input v-model="sshForm.password" type="password" show-password></el-input>
            </el-form-item>
            <el-form-item label="私钥" v-if="sshForm.logintype === 1">
              <el-input type="textarea" v-model="sshForm.privateKey" :rows="4" placeholder="Paste Private Key Content"></el-input>
            </el-form-item>
            <el-form-item label="密钥密码" v-if="sshForm.logintype === 1">
              <el-input v-model="sshForm.passphrase" type="password" placeholder="Passphrase (optional)"></el-input>
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane label="命令列表" name="commands">
            <div style="margin-top: 20px;">
            <div v-for="(cmd, index) in commandList" :key="index" class="cmd-row">
               <el-input v-model="commandList[index]" placeholder="请输入命令" style="width: 80%"></el-input>
               <el-button type="danger" icon="el-icon-delete" circle @click="removeCommand(index)"></el-button>
            </div>
            <el-button type="dashed" style="width: 100%; margin-top: 10px" @click="addCommand">+ 添加命令步骤</el-button>
            <div class="tip">命令将按顺序执行，前一条失败将终止后续执行</div>
          </el-tab-pane>
        </el-tabs>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="saveJob" :loading="saving">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import request from '@/utils/request'

export default {
  name: 'CronList',
  data() {
    return {
      jobs: [],
      loading: false,
      saving: false,
      dialogVisible: false,
      activeTab: 'basic',
      form: {
        ID: 0,
        Name: '',
        CronExpr: '20 07 * * *',
        Status: 1,
        RandomDelay: 20,
        HostInfo: '',
        Commands: ''
      },
      // SSH 临时表单
      sshForm: {
        hostname: '',
        port: 22,
        username: 'root',
        password: '',
        logintype: 0,
        privateKey: '',
        passphrase: ''
      },
      commandList: [''],
      rules: {
        Name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
        CronExpr: [{ required: true, message: '请输入Cron表达式', trigger: 'blur' }]
      }
    }
  },
  computed: {
    dialogTitle() {
      return this.form.ID ? '编辑任务' : '新建任务'
    }
  },
  created() {
    this.fetchJobs()
  },
  methods: {
    async fetchJobs() {
      this.loading = true
      try {
        const res = await request.get('/cron/list')
        if (res.code === 200) {
          this.jobs = res.data
        }
      } catch (e) {
        console.error(e)
      } finally {
        this.loading = false
      }
    },
    handleAdd() {
      this.resetForm()
      this.dialogVisible = true
    },
    handleEdit(row) {
      this.resetForm()
      this.form = { ...row }
      // Parse HostInfo
      try {
        if (row.HostInfo) {
          const jsonStr = atob(row.HostInfo)
          this.sshForm = JSON.parse(jsonStr)
        }
      } catch (e) {
        console.error("Failed to parse host info", e)
      }
      // Parse Commands
      try {
        if (row.Commands) {
          this.commandList = JSON.parse(row.Commands)
        }
      } catch (e) {
        this.commandList = ['']
      }
      this.dialogVisible = true
    },
    async handleDelete(row) {
      try {
        await this.$confirm('确认删除该任务?', '提示', { type: 'warning' })
        const res = await request.post(`/cron/delete/${row.ID}`)
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.fetchJobs()
        }
      } catch (e) {
        // cancel
      }
    },
    async handleRun(row) {
      try {
        const res = await request.post(`/cron/run/${row.ID}`)
        if (res.code === 200) {
          this.$message.success('已触发执行')
          this.fetchJobs()
        }
      } catch (e) {
        console.error(e)
      }
    },
    async saveJob() {
      this.$refs.form.validate(async valid => {
        if (valid) {
          // Prepare data
          if (!this.sshForm.hostname) {
             this.$message.error('请输入SSH主机地址')
             return
          }
          // Encode SSH Info
          const sshJson = JSON.stringify(this.sshForm)
          this.form.HostInfo = btoa(sshJson)
          
          // Encode Commands
          const validCmds = this.commandList.filter(c => c.trim() !== '')
          if (validCmds.length === 0) {
             this.$message.error('至少输入一条命令')
             return
          }
          this.form.Commands = JSON.stringify(validCmds)

          this.saving = true
          try {
            const res = await request.post('/cron/save', this.form)
            if (res.code === 200) {
              this.$message.success('保存成功')
              this.dialogVisible = false
              this.fetchJobs()
            }
          } catch (e) {
            console.error(e)
          } finally {
            this.saving = false
          }
        }
      })
    },
    resetForm() {
      this.activeTab = 'basic'
      this.form = { ID: 0, Name: '', CronExpr: '20 07 * * *', Status: 1, RandomDelay: 20, HostInfo: '', Commands: '' }
      this.sshForm = { hostname: '', port: 22, username: 'root', password: '', logintype: 0, privateKey: '', passphrase: '' }
      this.commandList = ['']
    },
    addCommand() {
      this.commandList.push('')
    },
    removeCommand(index) {
      this.commandList.splice(index, 1)
    },
    formatTime(time) {
      if (!time) return '从未'
      return new Date(time).toLocaleString()
    },
    getResultType(result) {
       if (result === '成功') return 'success'
       if (result === '失败') return 'danger'
       return 'info'
    }
  }
}
</script>

<style scoped lang="scss">
.cron-list-container {
  padding: 20px;
}
.tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}
.cmd-row {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  button {
    margin-left: 10px;
  }
}
.log-preview {
  max-height: 300px;
  overflow: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
  background: #333;
  color: #eee;
  padding: 10px;
  border-radius: 4px;
}
</style>
