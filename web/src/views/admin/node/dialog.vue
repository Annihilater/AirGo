<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="769px" destroy-on-close align-center>
    <el-divider content-position="left">节点参数</el-divider>
    <el-form :model="dialogData.nodeInfo" label-width="120px">
      <el-form-item label="node_type">
        <el-radio-group v-model="dialogData.nodeInfo.node_type">
          <el-radio label="vmess">vmess</el-radio>
          <el-radio label="vless">vless</el-radio>
          <el-radio label="trojan">trojan</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="remarks">
        <el-input v-model="dialogData.nodeInfo.remarks"/>
      </el-form-item>

      <!--      <el-form-item label="uuid">-->
      <!--        <el-input v-model="dialogData.nodeInfo.uuid"/>-->
      <!--      </el-form-item>-->
      <el-form-item label="address">
        <el-input v-model="dialogData.nodeInfo.address"/>
      </el-form-item>
      <el-form-item label="port">
        <el-input v-model.number="dialogData.nodeInfo.port"/>
      </el-form-item>

      <el-form-item label="scy" v-if="dialogData.nodeInfo.node_type==='vmess'">
        <!--        <el-input v-model="dialogData.nodeInfo.scy"/>-->
        <el-radio-group v-model="dialogData.nodeInfo.scy">
          <el-radio label="auto">auto</el-radio>
          <el-radio label="none">none</el-radio>
          <el-radio label="chacha20-poly1305">chacha20-poly1305</el-radio>
          <el-radio label="aes-128-gcm">aes-128-gcm</el-radio>
          <el-radio label="zero">zero</el-radio>
        </el-radio-group>

      </el-form-item>
      <el-form-item label="aid" v-if="dialogData.nodeInfo.node_type==='vmess'">
        <el-input v-model="dialogData.nodeInfo.aid"/>
      </el-form-item>


      <el-form-item label="network">
        <el-radio-group v-model="dialogData.nodeInfo.network">
          <el-radio label="tcp">tcp</el-radio>
          <el-radio label="kcp">kcp</el-radio>
          <el-radio label="ws">ws</el-radio>
          <el-radio label="h2">h2</el-radio>
          <el-radio label="quic">quic</el-radio>
          <el-radio label="grpc">grpc</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="type"
                    v-if="dialogData.nodeInfo.network==='tcp' || dialogData.nodeInfo.network==='kcp' || dialogData.nodeInfo.network==='quic'">
        <el-radio-group v-model="dialogData.nodeInfo.type">
          <el-radio label="none"
                    v-if="dialogData.nodeInfo.network==='tcp' || dialogData.nodeInfo.network==='kcp' || dialogData.nodeInfo.network=='quic'">
            none
          </el-radio>
          <el-radio label="http" v-if="dialogData.nodeInfo.network==='tcp'">http</el-radio>
          <el-radio label="srtp"
                    v-if="dialogData.nodeInfo.network==='kcp' || dialogData.nodeInfo.network=='quic'">
            srtp
          </el-radio>
          <el-radio label="utp"
                    v-if="dialogData.nodeInfo.network==='kcp' || dialogData.nodeInfo.network=='quic'">
            utp
          </el-radio>
          <el-radio label="wechat-video"
                    v-if="dialogData.nodeInfo.network==='kcp' || dialogData.nodeInfo.network=='quic'">
            wechat-video
          </el-radio>
          <el-radio label="dtls"
                    v-if="dialogData.nodeInfo.network==='kcp' || dialogData.nodeInfo.network=='quic'">
            dtls
          </el-radio>
          <el-radio label="wireguard"
                    v-if="dialogData.nodeInfo.network==='kcp' || dialogData.nodeInfo.network=='quic'">
            wireguard
          </el-radio>
        </el-radio-group>
      </el-form-item>


      <el-form-item label="host">
        <el-input v-model="dialogData.nodeInfo.host"/>
      </el-form-item>
      <el-form-item label="path">
        <el-input v-model="dialogData.nodeInfo.path"/>
      </el-form-item>
      <el-form-item label="mode" v-if="dialogData.nodeInfo.network==='grpc'">
        <el-radio-group v-model="dialogData.nodeInfo.mode">
          <el-radio label="gun">gun</el-radio>
          <el-radio label="multi">multi</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="security">
        <el-radio-group v-model="dialogData.nodeInfo.security">
          <el-radio label="">none</el-radio>
          <el-radio label="tls">tls</el-radio>
          <el-radio label="reality">reality</el-radio>
        </el-radio-group>

      </el-form-item>
      <el-form-item label="sni" v-if="dialogData.nodeInfo.security!==''">
        <el-input v-model="dialogData.nodeInfo.sni"/>
      </el-form-item>
      <el-form-item label="fp" v-if="dialogData.nodeInfo.security!==''">
        <el-input v-model="dialogData.nodeInfo.fp"/>
      </el-form-item>
      <el-form-item label="alpn" v-if="dialogData.nodeInfo.security==='tls'">
        <el-input v-model="dialogData.nodeInfo.alpn"/>
      </el-form-item>
      <el-form-item label="allowInsecure" v-if="dialogData.nodeInfo.security==='tls'">
        <el-switch
            size="small"
            v-model="dialogData.nodeInfo.allowInsecure"
            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
        />
      </el-form-item>
      <el-form-item label="pbk" v-if="dialogData.nodeInfo.security==='reality'">
        <el-input v-model="dialogData.nodeInfo.pbk"/>
      </el-form-item>
      <el-form-item label="sid" v-if="dialogData.nodeInfo.security==='reality'">
        <el-input v-model="dialogData.nodeInfo.sid"/>
      </el-form-item>
      <el-form-item label="spx" v-if="dialogData.nodeInfo.security==='reality'">
        <el-input v-model="dialogData.nodeInfo.spx"/>
      </el-form-item>
    </el-form>


    <el-divider content-position="left">其他参数</el-divider>
    <el-form :model="dialogData.nodeInfo" label-width="120px">
      <el-form-item label="是否启用">
        <el-switch
            size="small"
            v-model="dialogData.nodeInfo.enabled"
            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
        />
      </el-form-item>
      <el-form-item label="节点限速">
        <el-input type="number" v-model.number="dialogData.nodeInfo.node_speedlimit"/>
      </el-form-item>
      <el-form-item label="节点倍率">
        <el-input type="number" v-model.number="dialogData.nodeInfo.traffic_rate"/>
      </el-form-item>
      <el-form-item label="启用中转">
        <el-switch
            size="small"
            v-model="dialogData.nodeInfo.enable_transfer"
            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
        />
      </el-form-item>
      <el-form-item label="中转ip" v-if="dialogData.nodeInfo.enable_transfer">
        <el-input v-model="dialogData.nodeInfo.transfer_address" placeholder=""/>
      </el-form-item>
      <el-form-item label="中转端口" v-if="dialogData.nodeInfo.enable_transfer">
        <el-input v-model="dialogData.nodeInfo.transfer_port" placeholder=""/>
      </el-form-item>
    </el-form>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="state.isShowDialog = false">取消</el-button>
                <el-button type="primary" @click="onSubmit">
                    确认
                </el-button>
            </span>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>

import {storeToRefs} from "pinia";
import {useNodeStore} from "/@/stores/nodeStore";
import {reactive} from "vue";

const nodeStore = useNodeStore()
const {dialogData} = storeToRefs(nodeStore)
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);
//定义参数
const state = reactive({
  type: "",
  title: "",
  isShowDialog: false,
})

// 打开弹窗
const openDialog = (type: string, row?: any) => {
  if (type == 'add') {
    dialogData.value.nodeInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    state.type = type
    state.title = "新建节点"
    state.isShowDialog = true
  } else {
    state.type = type
    state.title = "修改节点"
    dialogData.value.nodeInfo = row  //将当前row写入pinia
    state.isShowDialog = true
  }
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
function onSubmit() {
  if (state.type === 'add') {
    //新建节点
    nodeStore.newNode(dialogData.value.nodeInfo)
    setTimeout(() => {
      emit('refresh');
    }, 500);
  } else {
    //更新节点
    nodeStore.updateNode(dialogData.value.nodeInfo)
    setTimeout(() => {
      emit('refresh');
    }, 500);
  }
  closeDialog()
}

// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});
</script>


<style scoped>
.dialog-footer button:first-child {
  margin-right: 10px;
}
</style>
  