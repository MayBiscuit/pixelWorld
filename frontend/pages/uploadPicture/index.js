// pages/uploadPicture/index.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    newWorldId:'',
    gridSize: 25,
    pixelSize:'27.2rpx',
    pixelData: [],
    imagePath: '',
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(options) {
    let pages = getCurrentPages();
    console.log(pages);
    const newWorldId = options.newWorldId;
    console.log('跳转携带的wid: '+newWorldId);
    this.setData({
      newWorldId: newWorldId,
    });
  },

  onGridSizeChange(e) {
    this.setData({ gridSize: e.detail.value });
    const pixelSize = (680 / this.data.gridSize) + 'rpx';
    this.setData({ pixelSize });
    console.log(pixelSize);
    if(this.data.pixelData.length > 0){
      this.uploadToServer();
    }
  },
  navigateBack() {
    wx.navigateBack();
  },

  uploadPicture(){
    console.log('uploadPicture被调用');
    wx.chooseImage({
      count: 1, 
      sizeType: ['original', 'compressed'], 
      sourceType: ['album', 'camera'], 
      success: (res) => {
        const tempFilePaths = res.tempFilePaths;
        console.log('选择的图片路径:', tempFilePaths[0]);
        this.setData({
          imagePath: tempFilePaths[0] 
        });
        this.uploadToServer();
      },
      fail: (err) => {
        console.error('选择图片失败:', err);
      }
    });
  },
  uploadToServer() {
    console.log('开始上传图片:', this.data.imagePath);
    // 获取gridSize参数
    const gridSize = this.data.gridSize;
    // 使用wx.uploadFile上传图片和参数
    wx.uploadFile({
      url: 'http://127.0.0.1:9000/world/picture/upload',
      filePath: this.data.imagePath, 
      name: 'file', 
      formData: {
        'gridSize': gridSize.toString() 
      },
      success: (uploadRes) => {
        console.log('图片上传成功:', uploadRes);
        // 上传成功后的操作
        wx.showToast({
          title: '上传成功',
          icon: 'success',
          duration: 2000
        });
        console.log('uploadRes: '+uploadRes)
        const resData = JSON.parse(uploadRes.data);
        const numRows = resData.length;
        const numCols = resData[0].length;
        console.log('二维数组的大小: 行数 =', numRows, '列数 =', numCols);
        this.setData({
          pixelData: resData
        });
      },
      fail: (uploadErr) => {
        console.error('图片上传失败:', uploadErr);
        // 上传失败后的操作
        wx.showToast({
          title: '上传失败',
          icon: 'none',
          duration: 2000
        });
      }
    });
  },
  goStudio() {
    // console.log('携带的 pixelData 类型:', typeof this.data.pixelData);
    const pixelDataString = JSON.stringify(this.data.pixelData);
    console.log("创建跳转携带gridSize： " + this.data.gridSize);
  
    // 准备请求体
    const requestBody = {
      background: pixelDataString,
      wid: this.data.newWorldId,
      wsize: this.data.gridSize
    };

    console.log("请求前gridSize "+this.data.gridSize);
  
    // 发起 PUT 请求
    wx.request({
      method: 'PUT',
      url: 'http://127.0.0.1:9000/world/picture/confirm',
      data: requestBody,
      success: (res) => {
        console.log('请求成功:', res);
        // 请求成功后跳转页面
        wx.navigateTo({
          url: `/pages/studio/index?pixelData=${encodeURIComponent(pixelDataString)}&newWorldId=${this.data.newWorldId}&gridSize=${this.data.gridSize}`,
        });
      },
      fail: (err) => {
        console.error('请求失败:', err);
        // 请求失败时可以在这里处理错误，例如显示错误提示
      }
    });
  }
})