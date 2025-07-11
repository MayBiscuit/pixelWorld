// pages/aiGenerate/index.js
Page({
  data: {
    newWorldId:'',
    gridSize: 25,
    inputValue: '',
    pixelData:[],
    pixelSize:'27.2rpx',
  },
  onLoad(options) {
    let pages = getCurrentPages();
    console.log(pages);
    const newWorldId = options.newWorldId;
    console.log("newWorldId: "+newWorldId);
    this.setData({
      newWorldId: newWorldId,
    });
  },
  navigateBack() {
    wx.navigateBack();
  },
  onInputChange(e) {
    this.setData({
      inputValue: e.detail.value
    });
  },
  onGridSizeChange(e) {
    this.setData({ gridSize: e.detail.value });
    const pixelSize = (680 / this.data.gridSize) + 'rpx';
    this.setData({ pixelSize });
    if(this.data.pixelData.length > 0){
      this.aiGenerate();
    }
  },
  aiGenerate() {
    // 获取输入框的内容
    const description = this.data.inputValue;

    if (!description) {
      wx.showToast({
        title: '请输入底纹描述',
        icon: 'none',
        duration: 2000
      });
      return;
    }
    // console.log(description)
    // console.log(this.data.gridSize)
    wx.showLoading({
      title: '生成中...',
      mask: true // 防止用户点击穿透
    });
    // 发起请求
    wx.request({
      url: 'http://127.0.0.1:9000/world/ai/upload', // 替换为你的实际接口地址
      method: 'GET',
      data: {
        description: description,
        gridSize:this.data.gridSize
      },
      success: (res) => {
        // console.log('请求成功:', res.data);
        // wx.showToast({
        //   title: '请求成功',
        //   icon: 'success',
        //   duration: 2000
        // });
        this.setData({
          pixelData: res.data
        });
      },
      fail: (err) => {
        console.error('请求失败:', err);
        wx.showToast({
          title: '请求失败',
          icon: 'none',
          duration: 2000
        });
      },
      complete: () => {
        // 无论成功还是失败，都隐藏加载提示
        wx.hideLoading();
      }
    });
  },
  goStudio(){
    // console.log('携带的 pixelData 类型:', typeof this.data.pixelData);
    if(this.data.pixelData.length==0){
      wx.showToast({
        title: '还未生成底纹',
        icon: 'none',
        duration: 2000
      });
    }else{
      const pixelDataString = JSON.stringify(this.data.pixelData);
      wx.navigateTo({
        url: `/pages/studio/index?pixelData=${encodeURIComponent(pixelDataString)}&newWorldId=${this.data.newWorldId}&gridSize=${this.data.gridSize}`,
      });
    }
  },
})