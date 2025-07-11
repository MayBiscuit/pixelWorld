// utils/auth.js
const API_BASE = 'http://127.0.0.1:9000';

export const checkLogin = () => {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${API_BASE}/api/check-auth`,
      method: 'GET',
      withCredentials: true,
      success(res) {
        if (res.data.isLoggedIn) {
          resolve(res.data.user);
        } else {
          resolve(null);
        }
      },
      fail(err) {
        reject(err);
      }
    });
  });
};

export const login = (userInfo) => {
  return new Promise((resolve, reject) => {
    wx.login({
      success: (loginRes) => {
        if (loginRes.code) {
          wx.request({
            url: `${API_BASE}/api/login`,
            method: 'POST',
            data: {
              code: loginRes.code,
              userInfo: userInfo
            },
            withCredentials: true,
            success(res) {
              if (res.data.success) {
                resolve(res.data.user);
              } else {
                reject(new Error('登录失败'));
              }
            },
            fail(err) {
              reject(err);
            }
          });
        } else {
          reject(new Error('获取code失败'));
        }
      },
      fail(err) {
        reject(err);
      }
    });
  });
};

export const logout = () => {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${API_BASE}/api/logout`,
      method: 'POST',
      withCredentials: true,
      success(res) {
        if (res.data.success) {
          resolve();
        } else {
          reject(new Error('登出失败'));
        }
      },
      fail(err) {
        reject(err);
      }
    });
  });
};