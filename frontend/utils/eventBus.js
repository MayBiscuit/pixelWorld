export default function createBus() { // 定义并导出一个默认的函数 createBus，用于创建一个事件总线实例
  return {
    events: {}, // 存储事件和对应的回调函数，初始为空
    on(event, callback) { // 向事件event添加回调函数callback
      if (!this.events[event]) this.events[event] = []; 
      this.events[event].push(callback); 
    },
    off(event, callback) { // 向事件event移除回调函数callback
      if (!this.events[event]) return;
      if (!callback) this.events[event] = [];
      else {
        const index = this.events[event].indexOf(callback);
        if (index !== -1) this.events[event].splice(index, 1);
      }
    },
    emit(event, ...args) { // 为事件event依次调用每个回调函数callback
      if (this.events[event]) this.events[event].forEach((callback) => callback(...args));
    },
  };
}
