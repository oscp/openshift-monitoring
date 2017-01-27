import {Injectable} from '@angular/core';
import {Observer, Observable, Subject} from 'rxjs/Rx';
import {NotificationsService} from "angular2-notifications";

@Injectable()
export class SocketService {
    public websocket: Subject<any>;

    constructor(private notificationService: NotificationsService) {
        this.connectToUI();
    }

    private reconnectWebsocket() {
        let that = this;
        this.notificationService.error("Error on websocket", "Error on websocket. Reconnecting...");
        setTimeout(
            () => {
                console.log('reconnecting websocket');
                that.websocket = undefined;
                that.connectToUI();
            }
            , 1000
        );
    }

    private connectToUI() {
        let that = this;
        let socket = new WebSocket('ws://localhost:8080/ui');
        let observable = Observable.create(
            (observer: Observer<MessageEvent>) => {
                socket.onmessage = observer.next.bind(observer);
                socket.onerror = () => {
                    that.reconnectWebsocket();
                }
                socket.onclose = () => {
                    that.reconnectWebsocket();
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
