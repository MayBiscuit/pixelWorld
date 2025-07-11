// pages/studio/index.js
Page({
  data: {
    worldId:'',
    worldName:'画室名称',
    worldDesc:'画室描述',
    description: '',
    currentColor:'#000000',
    gridSize:25,
    pixelSize:'27.2rpx',
    showColorPickerFlag: false,
    showStatisticPopup: false,
    topColors: [],
    pixelData: [],
    popupVisible: false,
    drawing: false, // 是否正在绘画
    lastRow: null, // 上一次触摸的行
    lastCol: null,  // 上一次触摸的列
    row1: [
      { type: 'icon', iconType: 'clear-formatting' },
      {'type': 'circle', 'color': '#E57373'},
      {'type': 'circle', 'color': '#F06292'},
      {'type': 'circle', 'color': '#BA68C8'},
      {'type': 'circle', 'color': '#9575CD'},
      {'type': 'circle', 'color': '#7986CB'},
      {'type': 'circle', 'color': '#64B5F6'},
      {'type': 'circle', 'color': '#4FC3F7'}
    ],
    row2: [
      { type: 'icon',  iconType: 'palette' },
      {'type': 'circle', 'color': '#4FC3F7'},
      {'type': 'circle', 'color': '#4DD0E1'},
      {'type': 'circle', 'color': '#4DB6AC'},
      {'type': 'circle', 'color': '#81C784'},
      {'type': 'circle', 'color': '#AED581'},
      {'type': 'circle', 'color': '#DCE775'},
      {'type': 'circle', 'color': '#FFF176'}
    ],
    row3: [
      { type: 'icon',  iconType: 'gift' },
      {'type': 'circle', 'color': '#FFD54F'},
      {'type': 'circle', 'color': '#FFB74D'},
      {'type': 'circle', 'color': '#FF8A65'},
      {'type': 'circle', 'color': '#A1887F'},
      {'type': 'circle', 'color': '#E0E0E0'},
      {'type': 'circle', 'color': '#90A4AE'},
      {'type': 'circle', 'color': '#000000'}
    ]
  },
  generatePixelData(gridSize) {
    console.log(`生成的数组大小: ${gridSize} x${gridSize}`); 
    const pixelData = Array.from({ length: gridSize }, () =>
    Array.from({ length: gridSize }, () => '#FFFFFF')
    );
    console.log('生成的像素数据:', pixelData); 
  return pixelData;
  },  
  onLoad(options) {
    let pages = getCurrentPages();
    console.log(pages);

    const oldWorldId = options.oldWorldId;
    console.log("oldWorldId: "+oldWorldId);
    if (oldWorldId) {// 从home页面跳转旧画室
      this.setData({
        worldId: oldWorldId,
      });
      wx.request({
        url: `http://127.0.0.1:9000/draw/thisworld/${oldWorldId}`, 
        method: 'GET',
        success: (res) => {
          if (res.statusCode === 200) {
            this.setData({
              gridSize: res.data.wsize,
              pixelData: res.data.wpicture,
              worldName:res.data.wname,
              worldDesc:res.data.wdesc,
              pixelSize:(680 / res.data.wsize) + 'rpx'
            });
          } else {
            console.error('请求失败:', res.statusCode);
          }
        },
        fail: (err) => {
          console.error('请求失败:', err);
        }
      });
    }else{ 
      const newWorldId = options.newWorldId;
      console.log("创建新画室： "+newWorldId);
      const gridSize = options.gridSize;
      const pixelSize = (680 / gridSize) + 'rpx';
      this.data.gridSize = gridSize;
      this.data.worldId = newWorldId;
      this.setData({ pixelSize });
  
      if(options.pixelData){ // 从其他模板跳转
        const pixelData = JSON.parse(decodeURIComponent(options.pixelData));
        this.setData({
          pixelData: pixelData
        });
      }else{
        this.setData({ // 从空白模板跳转
          pixelData: this.generatePixelData(this.data.gridSize)
        });
      }

      wx.request({
        url: `http://127.0.0.1:9000/draw/thisworld/${newWorldId}`, 
        method: 'GET',
        success: (res) => {
          if (res.statusCode === 200) {
            this.setData({
              worldName:res.data.wname,
              worldDesc:res.data.wdesc,
            });
          } else {
            console.error('请求失败:', res.statusCode);
          }
        },
        fail: (err) => {
          console.error('请求失败:', err);
        }
      });
  
      const wid = Number(this.data.worldId); 
      const background = JSON.stringify(this.data.pixelData); 
      wx.request({
        url: 'http://127.0.0.1:9000/draw/draw',
        method: 'PUT',
        header: {
          'Content-Type': 'application/json'
        },
        data: {
          background: background,
          wid: wid,
          wsize: parseInt(this.data.gridSize, 10)
        },
        success: (res) => {
          console.log('PUT 请求成功', res.data);
        },
        fail: (err) => {
          console.error('PUT 请求失败', err);
        }
      });
    }
  },
  editWorldName() {
    wx.showModal({
      placeholderText: ' 请输入新的名称',
      editable: true,
      success: (res) => {
        if (res.confirm) {
          const inputValue = res.content;
          console.log("inputValue: " + inputValue);
          if (!inputValue || inputValue.trim() === '') {
            wx.showToast({
              title: '请输入有效内容',
              icon: 'none'
            });
            return;
          }
          const wid = Number(this.data.worldId);
          const wname = inputValue;
          wx.request({
            url: 'http://127.0.0.1:9000/draw/modifyname',
            method: 'PUT',
            header: {
              'Content-Type': 'application/json'
            },
            data: {
              wname: wname,
              wid: wid
            },
            success: (res) => {
              this.setData({ worldName: wname });
            },
            fail: (err) => {
              console.error('PUT 请求失败', err);
            }
          });
        }
      },
      fail: (error) => {
        console.error('弹窗失败:', error);
      }
    });
  },  
  editWorldDesc(){
    console.log("调用编辑描述");
    wx.showModal({
      placeholderText: ' 请输入新的描述',
      editable: true,
      success: (res) => {
        if (res.confirm) {
          const inputValue = res.content;
          console.log("inputValue: " + inputValue);
          if (!inputValue || inputValue.trim() === '') {
            wx.showToast({
              title: '请输入有效内容',
              icon: 'none'
            });
            return;
          }
          const wid = Number(this.data.worldId);
          const wdesc = inputValue;
          wx.request({
            url: 'http://127.0.0.1:9000/draw/modifydesc',
            method: 'PUT',
            header: {
              'Content-Type': 'application/json'
            },
            data: {
              wdesc: wdesc,
              wid: wid
            },
            success: (res) => {
              this.setData({ worldDesc: wdesc });
            },
            fail: (err) => {
              console.error('PUT 请求失败', err);
            }
          });
        }
      },
      fail: (error) => {
        console.error('弹窗失败:', error);
      }
    });
  },
  changeWorldStatus: function() {
    const wid = this.data.worldId;
  
    // 弹出确认对话框
    wx.showModal({
      title: '确认操作',
      content: '是否要切换画室状态？',
      showCancel: true,
      success: (modalRes) => {
        if (modalRes.confirm) {
          console.log('用户点击了“确定”按钮');
          // 构造请求 URL
          const url = `http://127.0.0.1:9000/draw/changeworldstatus/${wid}`;
  
          // 发送 PUT 请求
          wx.request({
            url: url,
            method: 'PUT',
            header: {
              'Content-Type': 'application/json'
            },
            success: (res) => {
              console.log('PUT 请求成功', res.data);
              if (res.statusCode === 200) {
                wx.showToast({
                  title: '状态更新成功',
                  icon: 'success',
                  duration: 2000
                });
              } else {
                wx.showToast({
                  title: '更新失败',
                  icon: 'none',
                  duration: 2000
                });
              }
            },
            fail: (err) => {
              console.error('PUT 请求失败', err);
              wx.showToast({
                title: '请求失败',
                icon: 'none',
                duration: 2000
              });
            }
          });
        } else if (modalRes.cancel) {
          console.log('用户点击了“取消”按钮');
          wx.showToast({
            title: '操作已取消',
            icon: 'none',
            duration: 2000
          });
        }
      }
    });
  },
  deleteWorld: function() {
    const wid = this.data.worldId;

    // 弹出确认对话框
    wx.showModal({
      title: '确认操作',
      content: '是否删除该画室？',
      showCancel: true,
      success: (modalRes) => {
        if (modalRes.confirm) {
          console.log('用户点击了“确定”按钮');
          // 构造请求 URL
          const url = `http://127.0.0.1:9000/draw/deleteworld/${wid}`;

          // 发送 DELETE 请求
          wx.request({
            url: url,
            method: 'DELETE',
            header: {
              'Content-Type': 'application/json'
            },
            success: (res) => {
              console.log('DELETE 请求成功', res.data);
              if (res.statusCode === 200) {
                wx.showToast({
                  title: '删除成功',
                  icon: 'success',
                  duration: 2000
                });
                // 确保在成功提示后跳转
                setTimeout(() => {
                  wx.switchTab({
                    url: '/pages/home/index',
                  });
                }, 2000); // 等待 2 秒后跳转
              } else {
                wx.showToast({
                  title: '删除失败',
                  icon: 'none',
                  duration: 2000
                });
              }
            },
            fail: (err) => {
              console.error('DELETE 请求失败', err);
              wx.showToast({
                title: '请求失败',
                icon: 'none',
                duration: 2000
              });
            }
          });
        } else if (modalRes.cancel) {
          console.log('用户点击了“取消”按钮');
          wx.showToast({
            title: '操作已取消',
            icon: 'none',
            duration: 2000
          });
        }
      }
    });
  },
  changeCurrentColor: function(e){
    const color = e.currentTarget.dataset.color; // 获取 item.color
    console.log("当前颜色：", color);
    this.data.currentColor=color;
  },
  changeWhiteColor: function(e){
    this.data.currentColor="#FFFFFF";
  },
  togglePixel: function(e) {
    const row = e.currentTarget.dataset.row;
    const col = e.currentTarget.dataset.col;
    const currentColor = this.data.currentColor;

    const newPixelData = [...this.data.pixelData];
    newPixelData[row][col] = currentColor;
    this.setData({ pixelData: newPixelData });

    const wid = Number(this.data.worldId); 
    const background = JSON.stringify(this.data.pixelData); 
    wx.request({
      url: 'http://127.0.0.1:9000/draw/draw',
      method: 'PUT',
      header: {
        'Content-Type': 'application/json'
      },
      data: {
        background: background,
        wid: wid,
        wsize: parseInt(this.data.gridSize, 10)
      },
      success: (res) => {
        console.log('PUT 请求成功', res.data);
      },
      fail: (err) => {
        console.error('PUT 请求失败', err);
      }
    });
  },
  startDrawing: function(e) {
    console.log(e.touches[0]);
    const touch = e.touches[0];
    const query = wx.createSelectorQuery();
    query.selectViewport().node().exec((res) => {
    const canvas = res[0].node;
    const ctx = canvas.getContext('2d');
    const pageX = touch.pageX;
    const pageY = touch.pageY;
    // 计算触摸点所在的像素行和列
    const row = Math.floor(pageY / this.data.pixelSize);
    const col = Math.floor(pageX / this.data.pixelSize);
    this.setData({
    drawing: true,
    lastRow: row,
    lastCol: col
    });
    this.drawPixel(row, col);
    });
    },
    
    // 绘制像素
    drawPixel: function(row, col) {
    const currentColor = this.data.currentColor;
    const newPixelData = [...this.data.pixelData];
    newPixelData[row][col] = currentColor;
    this.setData({ pixelData: newPixelData });
    this.saveDrawing();
    },
    
    // 绘制中
    draw: function(e) {
    if (!this.data.drawing) return;
    
    const touch = e.touches[0];
    const pageX = touch.pageX;
    const pageY = touch.pageY;
    // 计算触摸点所在的像素行和列
    const row = Math.floor(pageY / this.data.pixelSize);
    const col = Math.floor(pageX / this.data.pixelSize);
    
    if (row === this.data.lastRow && col === this.data.lastCol) {
    return; // 避免重复绘制同一个像素
    }
    
    this.setData({
    lastRow: row,
    lastCol: col
    });
    this.drawPixel(row, col);
    },
    
    // 结束绘画
    endDrawing: function() {
    this.setData({
    drawing: false,
    lastRow: null,
    lastCol: null
    });
    },
    
    // 保存绘画
    saveDrawing: function() {
    const wid = Number(this.data.worldId);
    const background = JSON.stringify(this.data.pixelData);
    
    wx.request({
    url: 'http://127.0.0.1:9000/draw/draw',
    method: 'PUT',
    header: {
    'Content-Type': 'application/json'
    },
    data: {
    background: background,
    wid: wid,
    wsize: parseInt(this.data.gridSize, 10)
    },
    success: (res) => {
    console.log('PUT 请求成功 '+ res.data);
    },
    fail: (err) => {
    console.error('PUT 请求失败 '+ err);
    }
    });
    },
  savePicture: function() {
    const ctx = wx.createCanvasContext('pixelCanvas');
    const { pixelData, pixelSize } = this.data;
    const pixelSizeNum = parseFloat(pixelSize);

    const width = pixelData[0].length * pixelSizeNum;
    const height = pixelData.length * pixelSizeNum;
    ctx.clearRect(0, 0, width, height);

    pixelData.forEach((row, rowIndex) => {
      row.forEach((color, colIndex) => {
        ctx.setFillStyle(color);
        ctx.fillRect(colIndex * pixelSizeNum, rowIndex * pixelSizeNum, pixelSizeNum, pixelSizeNum);
      });
    });
    
    ctx.draw(false, () => {
      setTimeout(() => {
        wx.canvasToTempFilePath({
          canvasId: 'pixelCanvas',
          success: (res) => {
            wx.saveImageToPhotosAlbum({
              filePath: res.tempFilePath,
              success: () => {
                wx.showToast({
                  title: '保存成功',
                  icon: 'success'
                });
              },
              fail: (err) => {
                console.error('保存到相册失败', err);
                wx.showToast({
                  title: '保存失败',
                  icon: 'none'
                });
              }
            });
          },
          fail: (err) => {
            console.error('生成图片失败', err);
            wx.showToast({
              title: '生成图片失败',
              icon: 'none'
            });
          }
        });
      }, 500); // 延迟100毫秒确保绘制完成
    });
  },
  navigateBack() {
    wx.navigateBack();
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
  // 显示调色盘
  showColorPicker: function () {
    this.setData({
      showColorPickerFlag: true
    });
  },
  // 颜色变化回调
  onColorChange(e) {
    console.log('颜色变化', e.detail);
    const { value } = e.detail;
    console.log("调色盘颜色值： ", value)
    this.setData({
      currentColor: value // 更新当前颜色
    });
  },
  // 关闭调色盘回调
  onColorPickerClose(e) {
    console.log('关闭调色盘', e.detail);
    this.setData({
      showColorPickerFlag: false // 关闭调色盘
    });
  },
  // 生成随机颜色
  randomChangeColor: function() {
    // 生成随机颜色
    const randomColor = this.generateRandomColor();
    this.setData({
      currentColor: randomColor
    });
    console.log("currentColor: "+this.data.currentColor)
  },
  // 生成随机十六进制颜色字符串
  generateRandomColor: function() {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
  },
  // 显示统计数据
  showCount() {
    // 统计颜色
    const colorCount = {};
    let totalColors = 0;

    this.data.pixelData.forEach(row => {
      row.forEach(color => {
        if (colorCount[color]) {
          colorCount[color]++;
        } else {
          colorCount[color] = 1;
        }
        totalColors++;
      });
    });

    // 转换为数组并排序
    const sortedColors = Object.keys(colorCount).map(color => ({
      color: color,
      count: colorCount[color]
    })).sort((a, b) => b.count - a.count);

    // 获取前五使用最多的颜色
    const topColors = sortedColors.slice(0, 5).map(item => {
      const percentage = ((item.count / totalColors) * 100).toFixed(2);
      console.log(`颜色: ${item.color}, 百分比: ${percentage}%`); // 打印百分比
      return {
        color: item.color,
        count: item.count,
        percentage: percentage
      };
    });

    // 显示弹窗
    this.setData({
      showStatisticPopup: true,
      topColors: topColors
    });
  },
  // 关闭统计数据弹窗
  closeStatisticPopup() {
    this.setData({
      showStatisticPopup: false
    });
  },
  goHome(){
    const wid = Number(this.data.worldId); 
    const background = JSON.stringify(this.data.pixelData); 
    wx.request({
      url: 'http://127.0.0.1:9000/draw/draw',
      method: 'PUT',
      header: {
        'Content-Type': 'application/json'
      },
      data: {
        background: background,
        wid: wid,
        wsize: parseInt(this.data.gridSize, 10)
      },
      success: (res) => {
        console.log('PUT 请求成功', res.data);
      },
      fail: (err) => {
        console.error('PUT 请求失败', err);
      }
    });
    wx.switchTab({
      url: '/pages/home/index',
    })
  }
})