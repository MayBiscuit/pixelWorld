// pages/release/index.js
Page({
  data: {
    popupVisible: false,
    gridSize: 25,
    name: '',       // 存储画室名称
    description: '', // 存储画室描述
    newWorldId:'', // 标识创建的新画室
    showButtonGroup: false
  },
  onLoad(options) {
    let pages = getCurrentPages();
    console.log(pages);
  },
  // 处理画室名称输入
  handleNameInput: function(e) {
    this.setData({ name: e.detail.value });
  },
  // 处理画室描述输入
  handleDescInput: function(e) {
    this.setData({ description: e.detail.value });
  },  
  // 创建画室
  createWorld: function() {
    const { name, description } = this.data;

    const app = getApp();
    const userid = app.globalData.openid;

    if (!name || !description) {
      wx.showToast({
        title: '请填写完整信息',
        icon: 'none'
      });
      return;
    }

    wx.request({
      url: 'http://127.0.0.1:9000/world/createworld', 
      method: 'POST',
      data: {
        userid: userid,
        wname: name,         
        wdesc: description
      },
      success: (res) => {
        if (res.statusCode === 200) {
          // wx.showToast({
          //   title: '创建成功',
          //   icon: 'success'
          // });
          
          const wid = res.data.world.wid;
          this.setData({
            newWorldId: wid  
          });
          this.setData({ showButtonGroup: true });
          // console.log('newWorldId: '+this.data.newWorldId)
        } else {
          wx.showToast({
            title: '创建失败',
            icon: 'none'
          });
        }
      },
      fail: (err) => {
        wx.showToast({
          title: '请求失败',
          icon: 'none'
        });
        console.error('请求失败:', err);
      }
    });
  },
  // 选择按钮
  goPreTemplate() {
    wx.navigateTo({
      url: `/pages/preTemplate/index?newWorldId=${this.data.newWorldId}`,
    });
  },
  goUploadPicture() {
    wx.navigateTo({
      url: `/pages/uploadPicture/index?newWorldId=${this.data.newWorldId}`,
    });
  },
  goAiGenerate() { 
    wx.navigateTo({
      url: `/pages/aiGenerate/index?newWorldId=${this.data.newWorldId}`,
    });
  },
  goEmptyTemplate() {
    this.setData({
      popupVisible: true
    });
  },
  closePopup() {
    this.setData({
      popupVisible: false
    });
  },
  onGridSizeChange(e) {
    this.setData({ gridSize: e.detail.value });
    console.log('修改后的gridSize: '+ this.data.gridSize)
  },
  // goStudio() { 
  //   if (this.data.newWorldId && this.data.gridSize) {
  //     wx.navigateTo({
  //       url: `/pages/studio/index?newWorldId=${this.data.newWorldId}&gridSize=${this.data.gridSize}`,
  //     });
  //   } else {
  //     console.error('newWorldId 或 gridSize 不存在');
  //   }
  // },
  goStudio() { 
    console.log('goStudio被调用');
    if (this.data.newWorldId && this.data.gridSize) {
      // 发送 PUT 请求
      console.log('发送前的格子数 '+this.data.gridSize)
      wx.request({
        url: 'http://127.0.0.1:9000/world/confirmEmpty',
        method: 'PUT',
        data: {
          wid: this.data.newWorldId,
          gridSize:this.data.gridSize
        },
        header: {
          'content-type': 'application/json' // 明确指定请求头
        },
        success: (res) => {
          // 请求成功后跳转页面
          wx.navigateTo({
            url: `/pages/studio/index?newWorldId=${this.data.newWorldId}&gridSize=${this.data.gridSize}`,
          });
        },
        fail: (err) => {
          console.error('请求失败:', err);
        }
      });
    } else {
      console.error('newWorldId或gridSize不存在');
    }
  },
});
