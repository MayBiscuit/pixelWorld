Component({
  options: {
    styleIsolation: 'shared',
  },
  properties: {
    navType: {
      type: String,
      value: 'title',
    },
    titleText: String,
  },
  methods: {
    searchTurn() {
      wx.navigateTo({
        url: `/pages/search/index`,
      });
    },
    onChangeValue(e){
      console.log("搜索值修改", e.detail.value);
      this.triggerEvent('change', { value: e.detail.value});
    }
  },
});
