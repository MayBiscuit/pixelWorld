// app.js 小程序入口文件，全局生命周期管理（初始化、监听事件...）
// app.json 全局配置（窗口样式、导航栏、tabBar、页面路径...）
// app.less 全局样式
import config from './config';
import Mock from './mock/index';
import createBus from './utils/eventBus';
import { connectSocket, fetchUnreadNum } from './mock/chat';

if (config.isMock) {
  Mock();
}

App({
  globalData: {
    openid: null, 
  },
  onLaunch() {
    // 小程序启动时执行
    const updateManager = wx.getUpdateManager();

    updateManager.onCheckForUpdate((res) => {
      // console.log(res.hasUpdate)
    });

    updateManager.onUpdateReady(() => {
      wx.showModal({
        title: '更新提示',
        content: '新版本已经准备好，是否重启应用？',
        success(res) {
          if (res.confirm) {
            updateManager.applyUpdate();
          }
        },
      });
    });

    this.getUnreadNum(); // 初始化未读消息数量
    this.connect(); // 建立 WebSocket 连接
  },
  globalData: { // 全局数据存储
    userInfo: null,
    unreadNum: 0, // 未读消息数量
    socket: null, // SocketTask 对象
  },

  /** 全局事件总线 */
  eventBus: createBus(),

  /** 初始化WebSocket */
  connect() {
    const socket = connectSocket();
    socket.onMessage((data) => { // 监听消息，若为未读消息则更新计数
      data = JSON.parse(data);
      if (data.type === 'message' && !data.data.message.read) this.setUnreadNum(this.globalData.unreadNum + 1);
    });
    this.globalData.socket = socket;
  },

  /** 获取未读消息数量 触发事件通知（如 TabBar 显示未读红点）*/
  getUnreadNum() {
    // fetchUnreadNum() 是一个异步函数，返回一个 Promise 对象
    // .then(({ data }) => { 是一个 Promise 的 then 方法，用于处理异步操作成功后的结果
    // { data } 是对象解构的语法，它从返回的结果对象中提取 data 属性
    //  ({ data }) => { 是一个箭头函数
    fetchUnreadNum().then(({ data }) => { 
      this.globalData.unreadNum = data;
      this.eventBus.emit('unread-num-change', data);
    });
  },

  /** 设置未读消息数量 触发事件通知（如 TabBar 显示未读红点）*/
  setUnreadNum(unreadNum) {
    this.globalData.unreadNum = unreadNum;
    this.eventBus.emit('unread-num-change', unreadNum);
  },
});
