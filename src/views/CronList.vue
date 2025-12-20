<template>
  <div class="cron-container">
    <div class="header">
      <h2>定时任务管理</h2>
      <el-button type="primary" size="small" icon="el-icon-plus" @click="openDialog()">创建任务</el-button>
    </div>

    <el-table :data="jobs" border stripe v-loading="loading">
      <el-table-column prop="Name" label="任务名称" width="180"></el-table-column>
      <el-table-column label="服务器" width="200">
        <template slot-scope="scope">
          {{ parseHost(scope.row.HostInfo) }}
        </template>
      </el-table-column>
      <el-table-column prop="CronExpr" label="Cron表达式" width="150">
          <template slot-scope="scope">
              <el-tag size="small">{{ scope.row.CronExpr }}</el-tag>
          </template>
      </el-table-column>
      <el-table-column label="上次执行" width="160">
        <template slot-scope="scope">
           {{ scope.row.LastRunTime ? new Date(scope.row.LastRunTime).toLocaleString() : '从未' }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template slot-scope="scope">
            <el-tag :type="scope.row.LastResult === '成功' ? 'success' : (scope.row.LastResult ? 'danger' : 'info')">
                {{ scope.row.LastResult || '待运行' }}
            </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button size="mini" type="text" @click="runJob(scope.row)">立即执行</el-button>
          <el-button size="mini" type="text" @click="showLog(scope.row)">日志</el-button>
          <el-button size="mini" type="text" style="color: #f56c6c;" @click="deleteJob(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 编辑弹窗 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="任务名称">
          <el-input v-model="form.Name" placeholder="例如：每日备份"></el-input>
        </el-form-item>
        <el-form-item label="Cron表达式">
          <el-input v-model="form.CronExpr" placeholder="例如：0 0 2 * * * （每天凌晨2点）"></el-input>
          <div class="tip">分 时 日 月 周</div>
        </el-form-item>
        <el-form-item label="服务器信息">
             <el-input type="textarea" :rows="3" v-model="form.HostInfo" placeholder="粘贴SSH连接页生成的Base64 Hash值"></el-input>
        </el-form-item>
        <el-form-item label="启用状态">
             <el-switch v-model="form.Status" active-text="启用" inactive-text="禁用"></el-switch>
        </el-form-item>
        
        <el-divider content-position="left">命令列表 (按顺序执行)</el-divider>
        <div v-for="(cmd, idx) in commands" :key="idx" class="cmd-row">
            <el-input v-model="commands[idx]" placeholder="输入Linux命令"></el-input>
            <el-button type="danger" icon="el-icon-delete" circle size="mini" @click="removeCmd(idx)"></el-button>
        </div>
        <el-button type="text" icon="el-icon-plus" @click="addCmd">添加下一条命令</el-button>
      </el-form>
      <div slot="footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="saveJob">保 存</el-button>
      </div>
    </el-dialog>
    
    <!-- 日志弹窗 -->
    <el-dialog title="执行日志" :visible.sync="logVisible">
        <pre class="log-pre">{{ currentLog }}</pre>
    </el-dialog>
  </div>
</template>

<script>
import request from '@/utils/request'

export default {
  data() {
    return {
      jobs: [],
      loading: false,
      dialogVisible: false,
      logVisible: false,
      currentLog: '',
      form: { ID: 0, Name: '', CronExpr: '', HostInfo: '', Status: true },
      commands: [''],
      dialogTitle: '新建任务'
    }
  },
  created() {
    this.fetchJobs()
  },
  methods: {
    async fetchJobs() {
      this.loading = true
      const res = await request.get('/cron/list')
      this.jobs = res.data
      this.loading = false
    },
    parseHost(base64Str) {
        try {
            const json = JSON.parse(atob(base64Str))
            return `${json.username}@${json.hostname}`
        } catch(e) { return 'Invalid Host Info' }
    },
    openDialog(row) {
        if (row) {
            this.form = { ...row }
            try { this.commands = JSON.parse(row.Commands) } catch(e) { this.commands = [''] }
            this.dialogTitle = '编辑任务'
        } else {
            this.form = { ID: 0, Name: '', CronExpr: '', HostInfo: '', Status: true }
            this.commands = ['']
            this.dialogTitle = '新建任务'
        }
        this.dialogVisible = true
    },
    addCmd() { this.commands.push('') },
    removeCmd(i) { this.commands.splice(i, 1) },
    async saveJob() {
        const payload = { ...this.form, Commands: JSON.stringify(this.commands) }
        await request.post('/cron/save', payload)
        this.dialogVisible = false
        this.fetchJobs()
        this.$message.success('保存成功')
    },
    async deleteJob(row) {
        await this.$confirm('确认删除该任务?')
        await request.post(`/cron/delete/${row.ID}`)
        this.fetchJobs()
    },
    async runJob(row) {
        await request.post(`/cron/run/${row.ID}`)
        this.$message.success('触发成功，请稍后查看日志')
        setTimeout(this.fetchJobs, 2000)
    },
    showLog(row) {
        this.currentLog = row.ErrorLog || '暂无日志'
        this.logVisible = true
    }
  }
}
</script>

<style scoped>
.cron-container { background: #fff; padding: 20px; border-radius: 8px; min-height: 500px; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.cmd-row { display: flex; gap: 10px; margin-bottom: 10px; }
.tip { font-size: 12px; color: #999; margin-top: 5px; }
.log-pre { background: #333; color: #fff; padding: 10px; border-radius: 4px; overflow: auto; max-height: 400px; }
</style>
