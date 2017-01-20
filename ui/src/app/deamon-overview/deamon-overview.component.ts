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
                <tr><th>Ip</th></tr>
            </thead>
            <tbody>
                <tr *ngFor="let d of deamons">
                    <td>{{d.Addr}}</td>
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
        this.socket.next({Type: SocketType.ALL_DEAMONS})
    }

    ngOnInit() {
        this.socket.subscribe(
            message => {
                let data = JSON.parse(message.data);

                console.log(data);

                switch (data.Type.Name) {
                    case SocketType.ALL_DEAMONS.Name:
                        this.deamons = data.Message;
                        console.log(this.deamons);
                        break;
                }

            }
        );
    }
}
