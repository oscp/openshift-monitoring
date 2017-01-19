import { Injectable } from '@angular/core';
import {Subject, Observer, Observable} from 'rxjs/Rx';

@Injectable()
export class SocketService {
  private existingSocket: Subject<any>;

  constructor() { }

  public createOrGetWebsocket(): Subject<MessageEvent> {
    if (!this.existingSocket) {
      let socket = new WebSocket('ws://localhost:8080/ui');
      let observable = Observable.create(
          (observer: Observer<MessageEvent>) => {
            socket.onmessage = observer.next.bind(observer);
            socket.onerror = observer.error.bind(observer);
            socket.onclose = observer.complete.bind(observer);
            return socket.close.bind(socket);
          }
      );
      let observer = {
        next: (data: Object) => {
          this.waitForSocketConnection(socket, () => {
            socket.send(JSON.stringify(data));
          });
        }
      };

      this.existingSocket = Subject.create(observer, observable);
    }
    return this.existingSocket;
  }

  private waitForSocketConnection(socket, callback) {
    setTimeout(
        function () {
          if (socket.readyState === 1) {
            console.log('Connection is made');
            if (callback != null) {
              callback();
            }
            return;

          } else {
            console.log('wait for connection...');
            this.waitForSocketConnection(socket, callback);
          }

        }, 5); // wait 5 milisecond for the connection...
  }
}
