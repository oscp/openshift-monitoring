import {Component, OnInit} from '@angular/core';
import {Subject} from "rxjs";
import {SocketService} from "../socket.service";
import {SocketType} from "../shared/socket.types";

@Component({
    selector: 'app-deamon-overview',
    template: `
        <h4>Connected Deamons</h4>
        <table class="table table-striped">
            <thead class="thead-inverse">
                <tr><th>Type</th>
                <th>IP</th>
                <th>Port</th>
                </tr>
            </thead>
            <tbody>
                <tr *ngFor="let d of deamons">
                    <td>{{d.DeamonType}}</td>
                    <td>{{d.Addr}}</td>
                    <td>{{d.Port}}</td>
                </tr>
            </tbody>
        </table>
    `
})
export class DeamonOverviewComponent implements OnInit {
    private socket: Subject<any>;
    private deamons: any;

    constructor(private socketService: SocketService) {
        this.socket = socketService.createOrGetWebsocket();
        this.socket.next({Type: SocketType.ALL_DEAMONS});
    }

    ngOnInit() {
        this.socket.subscribe(
            message => {
                let data = JSON.parse(message.data);

                console.log(data.Type.Name, data.Message);

                switch (data.Type.Name) {
                    case SocketType.ALL_DEAMONS.Name:
                        if (data.Message) {
                            this.deamons = [];

                            Object.keys(data.Message).forEach(k => {
                                this.deamons.push(data.Message[k]);
                            })
                        }
                        break;
                    case SocketType.NEW_DEAMON.Name:
                        this.deamons.push(data.Message);
                        break;
                    case SocketType.DEAMON_LEFT.Name:
                        this.socket.next({Type: SocketType.ALL_DEAMONS});
                }
            }
        );
    }
}
