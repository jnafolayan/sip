class AppEvent {
  constructor(name) {
    this.name = name;
    this.subscribers = [];
  }

  subscribe(callback) {
    this.subscribers.push(callback);
  }

  unsubscribe(callback) {
    this.subscribers = this.subscribers.filter(cb => cb != callback);
  }

  fire(data) {
    this.subscribers.forEach(cb => cb(data));
    console.log(`Event fired: ${this.name}`, data);
  }
}