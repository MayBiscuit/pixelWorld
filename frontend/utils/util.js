// 将数字格式化为两位数字符串，如果数字是一位数，则在前面补零
const formatNumber = (n) => { 
  n = n.toString();
  return n[1] ? n : `0${n}`;
};

// 将日期对象格式化为 YYYY/MM/DD HH:MM:SS 的字符串
const formatTime = (date) => {
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  const hour = date.getHours();
  const minute = date.getMinutes();
  const second = date.getSeconds();

  return `${[year, month, day].map(formatNumber).join('/')} ${[hour, minute, second].map(formatNumber).join(':')}`;
};

// 复制到本地临时路径，方便预览
const getLocalUrl = (path, name) => {
  const fs = wx.getFileSystemManager();
  const tempFileName = `${wx.env.USER_DATA_PATH}/${name}`;
  fs.copyFileSync(path, tempFileName);
  return tempFileName;
};

module.exports = {
  formatTime,
  getLocalUrl,
};
