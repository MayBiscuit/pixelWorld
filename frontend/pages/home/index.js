import Message from 'tdesign-miniprogram/message/index';
import request from '~/api/request';

const app = getApp()

Page({
  data: {
    enable: false,
    swiperList: [],
    edcardInfo: [],
    ingcardInfo: [],
  },
  onLoad(option) {
    this.checkLogin();
  },
  onShow(){
    if (app.globalData.openid) {
      const userid = app.globalData.openid;
      this.getEdWorld(userid);
      this.getIngWorld(userid);
    }
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
  checkLogin: function() {
    if (app.globalData.openid) {
      return;
    } else{
      wx.showToast({
        title: '还未登录，请先登录',
        icon: 'none',
        duration: 2000
      });
      
      setTimeout(() => {
        wx.switchTab({
          url: '/pages/my/index'
        });
      }, 2000);
    }
  },
  getEdWorld: function(userid) {
    const edurl = `http://127.0.0.1:9000/home/edworld/${userid}`;
    wx.request({
      url: edurl,
      method: 'GET',
      header: {
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          const worldData = res.data;
          const downloadPromises = worldData.map(item => this.downloadImage(item.wid));
          console.log("downloadPromises "+downloadPromises);

          Promise.all(downloadPromises).then(tempFilePaths => {
            const transformedData = worldData.map((item, index) => ({
              worldId: item.wid,
              url: tempFilePaths[index], // 使用下载的临时文件路径
              name: item.wname,
              desc: item.wdesc
            }));
            this.setData({
              edcardInfo: transformedData
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
        console.error('GET 请求失败', err);
      }
    });
  },
  getIngWorld: function(userid) {
    const ingurl = `http://127.0.0.1:9000/home/ingworld/${userid}`;
    wx.request({
      url: ingurl,
      method: 'GET',
      header: {
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          const worldData = res.data;
      const downloadPromises = worldData.map(item => this.downloadImage(item.wid));
      console.log("downloadPromises "+downloadPromises);

      Promise.all(downloadPromises).then(tempFilePaths => {
        const transformedData = worldData.map((item, index) => ({
          worldId: item.wid,
          url: tempFilePaths[index], // 使用下载的临时文件路径
          name: item.wname,
          desc: item.wdesc
        }));
        this.setData({
          ingcardInfo: transformedData
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
        console.error('GET 请求失败', err);
      }
    });
  },
  downloadImage: function(wid) {
    return new Promise((resolve, reject) => {
      const imageUrl = 'http://127.0.0.1:9000/draw/save?wid=' + wid;
      wx.request({
        url: imageUrl,
        method: 'GET',
        responseType: 'arraybuffer',
        success: (res) => {
          const tempFilePath = wx.env.USER_DATA_PATH + '/' + wid + '.png';
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
  onRefresh() {
    this.refresh();
  },
  async refresh() {
    console.log("刷新")
    if (app.globalData.openid) {
      const userid = app.globalData.openid;
      this.getEdWorld(userid);
      this.getIngWorld(userid);
    }
  },
  navigateToWorld: function(e){
    const data = e.currentTarget.dataset.worldid;//这里I不能大写。。。
    console.log("跳转前worldId: "+data);
    const url = `/pages/studio/index?oldWorldId=${encodeURIComponent(data)}`; 
    wx.navigateTo({
      url: url 
    });
  },
  goRelease() {
    wx.navigateTo({
      url: '/pages/release/index',
    });
  },
  searchWorld(){
    wx.showModal({
      placeholderText: ' 请输入搜索内容',
      editable: true, 
      success: (res) => {
        if (res.confirm) {
          const inputValue = res.content; 
          console.log("inputValue: "+ inputValue)
          if (!inputValue || inputValue.trim() === '') {
            wx.showToast({
              title: '请输入有效内容',
              icon: 'none'
            });
            return;
          }
  
          const openid = getApp().globalData.openid;
          if (!openid) {
            wx.showToast({
              title: '未获取到用户ID',
              icon: 'none'
            });
            return;
          }
          this.searchedWorld(openid, inputValue);
          this.searchingWorld(openid, inputValue);
        }
      },
      fail: (error) => {
        console.error('弹窗失败:', error);
      }
    });
  },
  searchedWorld: function(openid, key) {
    const url = `http://127.0.0.1:9000/home/searchedworld/${openid}`;
    wx.request({
      url: url,
      method: 'GET',
      data: {
        key: key
      },
      success: (response) => {
        // console.log('搜索已结束世界请求成功:', response.data);
        const worldData = response.data;
      const downloadPromises = worldData.map(item => this.downloadImage(item.wid));
      // console.log("downloadPromises "+downloadPromises);

      Promise.all(downloadPromises).then(tempFilePaths => {
        const transformedData = worldData.map((item, index) => ({
          worldId: item.wid,
          url: tempFilePaths[index], // 使用下载的临时文件路径
          name: item.wname,
          desc: item.wdesc
        }));
        this.setData({
          edcardInfo: transformedData
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
        console.error('搜索已结束世界请求失败:', error);
        wx.showToast({
          title: '搜索失败',
          icon: 'none'
        });
      }
    });
  },
  searchingWorld: function(openid, key) {
    const ingurl = `http://127.0.0.1:9000/home/searchingworld/${openid}`;
    wx.request({
      url: ingurl,
      method: 'GET',
      data: {
        key: key
      },
      success: (response) => {
        console.log('搜索进行中世界请求成功:', response.data);
        const worldData = response.data;
      const downloadPromises = worldData.map(item => this.downloadImage(item.wid));
      // console.log("downloadPromises "+downloadPromises);

      Promise.all(downloadPromises).then(tempFilePaths => {
        const transformedData = worldData.map((item, index) => ({
          worldId: item.wid,
          url: tempFilePaths[index], // 使用下载的临时文件路径
          name: item.wname,
          desc: item.wdesc
        }));
        this.setData({
          ingcardInfo: transformedData
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
        console.error('搜索进行中世界请求失败:', error);
        wx.showToast({
          title: '搜索失败',
          icon: 'none'
        });
      }
    });
  }
});
