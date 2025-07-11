// pages/preTemplate/index.js
Page({
  data: {
    newWorldId:'',
    popupVisible: false,
    gridSize: 25,
    backgroundInfo: [],
    gridSize: 25,
    pixelData:[],
    pixelSize:'27.2rpx',
    searchValue: ''
  },
  onLoad(options) {
    let pages = getCurrentPages();
    console.log(pages);
    const newWorldId = options.newWorldId;
    this.setData({
      newWorldId: newWorldId,
    });

    wx.request({
      url: 'http://127.0.0.1:9000/world/template/all', // 替换为你的API域名
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200) {
          const worldData = res.data;
      const downloadPromises = worldData.map(item => this.downloadImage(item.bid));
      // console.log("downloadPromises "+downloadPromises);

      Promise.all(downloadPromises).then(tempFilePaths => {
        const transformedData = worldData.map((item, index) => ({
          tid: item.bid,
          picture:item.bpicture,
          url: tempFilePaths[index], // 使用下载的临时文件路径
          name: item.bname,
          type: item.btype,
        }));
        this.setData({
          backgroundInfo: transformedData
        });
      }).catch(error => {
          console.error('下载图片失败:', error);
          wx.showToast({
            title: '图片下载失败',
            icon: 'none'
          });
        });
      }},
      fail: (err) => {
        // 请求失败，打印错误信息
        console.error('请求失败', err);
      },
    });
  },
  downloadImage: function(tid) {
    return new Promise((resolve, reject) => {
      const imageUrl = 'http://127.0.0.1:9000/world/template/save?tid=' + tid;
      wx.request({
        url: imageUrl,
        method: 'GET',
        responseType: 'arraybuffer',
        success: (res) => {
          const tempFilePath = wx.env.USER_DATA_PATH + '/' + tid + '.png';
          const fs = wx.getFileSystemManager();
          fs.writeFileSync(tempFilePath, res.data, 'binary');
          resolve(tempFilePath);
        },
        fail: (error) => {
          reject(error);
        }
      });
    });
  },
  onUnload: function() {
    // 小程序页面卸载时删除临时文件
    const fs = wx.getFileSystemManager();
    const tempFilePath = this.data.tempFilePath;
    if (tempFilePath) {
      fs.unlink({
        filePath: tempFilePath,
        success: (res) => {
          console.log('临时文件删除成功');
        },
        fail: (err) => {
          console.log('临时文件删除失败', err);
        }
      });
    }
  },
  onImageTap(e) {
    const bpicture = e.currentTarget.dataset.bpicture;
    this.setData({
      pixelData: bpicture,
    });
    // console.log("pixelData updated to:", this.data.pixelData);
    const gridSize = bpicture.length;
    const pixelSize = (680 / gridSize) + 'rpx';
    this.setData({
      gridSize: gridSize,
      pixelSize: pixelSize
    });
    // console.log("gridSize updated to:", this.data.gridSize);
    // console.log("pixelSize updated to:", this.data.pixelSize);
  },
  navigateBack() {
    wx.navigateBack();
  },
  goStudio(){
    // console.log('携带的 pixelData 类型:', typeof this.data.pixelData);
    if(this.data.pixelData.length==0){
      wx.showToast({
        title: '还未选择模板',
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
  openPopup() {
    this.setData({
      popupVisible: true
    });
  },
  closePopup() {
    this.setData({
      popupVisible: false
    });
  },
  onTabTap(e) {
    const key = e.currentTarget.dataset.key; // 获取按钮的文本值
    console.log("搜索关键词:", key);

    // 发起GET请求
    wx.request({
      url: 'http://127.0.0.1:9000/world/template/alltype', // 替换为你的API域名
      method: 'GET',
      data: {
        type: key, // 将按钮的文本值作为query参数key的值
      },
      success: (res) => {
        if (res.statusCode === 200) {
          // 请求成功，处理返回的数据
          console.log("搜索结果:", res.data);
          // 你可以根据需要更新页面的data对象
          const templateData = res.data;
          const downloadPromises = templateData.map(item => this.downloadImage(item.bid));
          // console.log("downloadPromises "+downloadPromises);
    
          Promise.all(downloadPromises).then(tempFilePaths => {
            const transformedData = templateData.map((item, index) => ({
              tid: item.bid,
              picture:item.bpicture,
              url: tempFilePaths[index], // 使用下载的临时文件路径
              name: item.bname,
              type: item.btype,
            }));
            this.setData({
              backgroundInfo: transformedData
            });
          }).catch(error => {
              console.error('下载图片失败:', error);
              wx.showToast({
                title: '图片下载失败',
                icon: 'none'
              });
            });
          }else {
          // 请求失败，打印错误信息
          console.error('请求失败', res);
        }
      },
      fail: (err) => {
        // 请求失败，打印错误信息
        console.error('请求失败', err);
      },
    });
  },
  onInput(e) {
    console.log("输入变化 "+e.detail.value);
    this.setData({
      searchValue: e.detail.value
    });
  },
  onSearch() {
    const { searchValue } = this.data;
    
    if (!searchValue.trim()) {
      wx.showToast({
        title: '请输入搜索内容',
        icon: 'none'
      });
      return;
    }

    // 发起请求
    wx.request({
      url: 'http://127.0.0.1:9000/world/template/search',
      method: 'GET',
      data: {
        key: searchValue
      },
      success: (res) => {
        console.log('搜索结果:', res.data);
        const templateData = res.data;
        const downloadPromises = templateData.map(item => this.downloadImage(item.bid));
        // console.log("downloadPromises "+downloadPromises);
  
        Promise.all(downloadPromises).then(tempFilePaths => {
          const transformedData = templateData.map((item, index) => ({
            tid: item.bid,
            picture:item.bpicture,
            url: tempFilePaths[index], // 使用下载的临时文件路径
            name: item.bname,
            type: item.btype,
          }));
          this.setData({
            backgroundInfo: transformedData
          });
        }).catch(error => {
            console.error('下载图片失败:', error);
            wx.showToast({
              title: '图片下载失败',
              icon: 'none'
            });
          });
      },
      fail: (error) => {
        console.error('请求失败:', error);
        wx.showToast({
          title: '搜索失败',
          icon: 'none'
        });
      }
    });
  }
})