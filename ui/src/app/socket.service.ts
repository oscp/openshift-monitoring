import {Injectable} from '@angular/core';
import {Observer, Observable, Subject} from 'rxjs/Rx';

@Injectable()
export class SocketService {
    public websocket: Subject<any>;

    constructor() {
        this.connectToUI();
    }

    private connectToUI() {
        let that = this;
        let socket = new WebSocket('ws://localhost:8080/ui');
        let observable = Observable.create(
            (observer: Observer<MessageEvent>) => {
                socket.onmessage = observer.next.bind(observer);
                socket.onerror = observer.error.bind(observer);
                socket.onclose = () => {
                    setTimeout(
                        () => {
                            console.log('reconnecting websocket');
                            that.websocket = undefined;
                            that.connectToUI();
                        }
                        , 10000
                    );
                };
                return socket.close.bind(socket);
            }
        ).share();

        let observer = {
            next: (data: Object) => {
                that.waitForSocketConnection(socket, () => {
                    socket.send(JSON.stringify(data));
                });
            }
        };

        this.websocket = Subject.create(observer, observable);
    }

    private waitForSocketConnection(socket, callback) {
        const that = this;
        setTimeout(
            function () {
                if (socket.readyState === 1) {
                    if (callback != null) {
                        callback();
                    }
                    return;

                } else {
                    console.log('wait for connection...');
                    that.waitForSocketConnection(socket, callback);
                }

            }, 5); // wait 5 milisecond for the connection...
    }
}
