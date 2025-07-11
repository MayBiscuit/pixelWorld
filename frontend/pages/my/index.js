import { checkLogin, login, logout } from '../../utils/auth';

Page({
  data: {
    userInfo: null,
    isLoggedIn: false,
    worldCount:0,
    // imgStr:'iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAIAAAACUFjqAAAAI0lEQVR4nGL6//8/AwMDLpIBjxyIxK+biYEgGJqGAwIAAP//mN1ooHbPSA0AAAAASUVORK5CYII=',
    // tempFilePath: ''
    // edcardInfo: [], // 存储最终要显示的数据
  },
  
  onLoad() {
    this.checkAuth();
    // const wid = '81';
    // const apiUrl = `http://127.0.0.1:9000/draw/save?wid=${wid}`;

    // // 发送 GET 请求
    // wx.request({
    //   url: apiUrl,
    //   method: 'GET',
    //   responseType: 'arraybuffer', // 接收二进制数据
    //   success: (res) => {
    //     if (res.statusCode === 200) {
    //       // 请求成功，获取二进制数据
    //       const arrayBuffer = res.data;
    //       // 生成唯一文件名
    //       const fs = wx.getFileSystemManager();
    //       const tempFilePath = `${wx.env.USER_DATA_PATH}/image_${new Date().getTime()}.png`;
    //       // 将二进制数据写入临时文件
    //       fs.writeFileSync(tempFilePath, arrayBuffer, 'binary');
    //       // 更新数据，显示图片
    //       this.setData({
    //         tempFilePath: tempFilePath
    //       });
    //     } else {
    //       // 请求失败
    //       wx.showToast({
    //         title: '下载图片失败',
    //         icon: 'none'
    //       });
    //     }
    //   },
    //   fail: (err) => {
    //     // 请求失败
    //     wx.showToast({
    //       title: '请求失败',
    //       icon: 'none'
    //     });
    //   }
    // });
  },

  // downloadImage: function(wid) {
  //   return new Promise((resolve, reject) => {
  //     const imageUrl = 'http://127.0.0.1:9000/draw/save?wid=' + wid;
  //     wx.request({
  //       url: imageUrl,
  //       method: 'GET',
  //       responseType: 'arraybuffer',
  //       success: (res) => {
  //         const tempFilePath = wx.env.USER_DATA_PATH + '/' + wid + '.png';
  //         const fs = wx.getFileSystemManager();
  //         fs.writeFileSync(tempFilePath, res.data, 'binary');
  //         resolve(tempFilePath);
  //       },
  //       fail: (error) => {
  //         reject(error);
  //       }
  //     });
  //   });
  // },

  // onUnload: function() {
  //   // 小程序页面卸载时删除临时文件
  //   const fs = wx.getFileSystemManager();
  //   const tempFilePath = this.data.tempFilePath;
  //   if (tempFilePath) {
  //     fs.unlink({
  //       filePath: tempFilePath,
  //       success: (res) => {
  //         console.log('临时文件删除成功');
  //       },
  //       fail: (err) => {
  //         console.log('临时文件删除失败', err);
  //       }
  //     });
  //   }
  // },

  async checkAuth() {
    try {
      const user = await checkLogin();
      if (user) {
        this.setData({
          userInfo: user,
          isLoggedIn: true
        });
      }
    } catch (err) {
      console.error('检查登录状态失败:', err);
    }
  },
  async handleLogin() {
    try {
      // 获取用户信息
      const userInfoRes = await new Promise((resolve, reject) => {
        wx.getUserProfile({
          desc: '用于完善会员资料',
          success: resolve,
          fail: reject
        });
      });
      
      // 调用登录
      const user = await login(userInfoRes.userInfo);
      
      this.setData({
        userInfo: user,
        isLoggedIn: true
      });

      const app = getApp();
      app.globalData.openid = user.openid;
      console.log(app.globalData.openid);
      
      wx.showToast({
        title: '登录成功',
        icon: 'success'
      });

      this.fetchWorldCount();
    } catch (err) {
      console.error('登录失败:', err);
      wx.showToast({
        title: '登录失败',
        icon: 'none'
      });
    }
    // const app = getApp();
    // const userid2 = app.globalData.openid; 
    // console.log("my userid "+userid2);
    // const url = `http://127.0.0.1:9000/home/edworld/${userid2}`;
    // wx.request({
    //   url: url,
    //   method: 'GET',
    //   header: {
    //     'Content-Type': 'application/json'
    //   },
    //   success: (response) => {
    //     console.log('my请求成功:', response.data);
    //     const worldData = response.data;
    //     const downloadPromises = worldData.map(item => this.downloadImage(item.wid));
    //     console.log("downloadPromises "+downloadPromises);

    //     Promise.all(downloadPromises).then(tempFilePaths => {
    //       const transformedData = worldData.map((item, index) => ({
    //         worldId: item.wid,
    //         url: tempFilePaths[index], // 使用下载的临时文件路径
    //         name: item.wname,
    //         desc: item.wdesc
    //       }));
    //       this.setData({
    //         edcardInfo: transformedData
    //       });
    //     }).catch(error => {
    //       console.error('下载图片失败:', error);
    //       wx.showToast({
    //         title: '图片下载失败',
    //         icon: 'none'
    //       });
    //     });
    //   },
    //   fail: (error) => {
    //     console.error('请求失败:', error);
    //     wx.showToast({
    //       title: '搜索失败',
    //       icon: 'none'
    //     });
    //   }
    // });
  },
  async handleLogout() {
    try {
      await logout();
      this.setData({
        userInfo: null,
        isLoggedIn: false
      });
      wx.showToast({
        title: '已退出登录',
        icon: 'success'
      });
    } catch (err) {
      console.error('登出失败:', err);
      wx.showToast({
        title: '登出失败',
        icon: 'none'
      });
    }
  },
  fetchWorldCount: function() {
    const app = getApp(); // 获取全局变量
    const userid = app.globalData.openid; // 获取全局变量中的 openid

    if (!userid) {
      console.error('用户未登录，无法获取 userid');
      return;
    }

    // 构造请求 URL
    const url = `http://127.0.0.1:9000/home/allworld/${userid}`;

    // 发送 GET 请求
    wx.request({
      url: url,
      method: 'GET',
      header: {
        'Content-Type': 'application/json'
      },
      success: (res) => {
        console.log('GET 请求成功', res.data);
        if (res.statusCode === 200) {
          const worldData = res.data; // 获取返回的数据
          const worldCount = worldData.length; // 计算数据中的元素个数
          this.setData({
            worldCount: worldCount // 更新 worldCount
          });
        } else {
          console.error('获取数据失败', res.data);
        }
      },
      fail: (err) => {
        console.error('GET 请求失败', err);
      }
    });
  }
});
