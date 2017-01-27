import {Component, OnInit} from '@angular/core';
import {SocketService} from "../socket.service";
import {SocketType} from "../shared/socket.types";

@Component({
    selector: 'app-checks',
    templateUrl: 'checks.component.html'
})
export class ChecksComponent implements OnInit {
    public checks = {};

    constructor(private socketService: SocketService) {
        this.getCurrentChecks();
    }

    ngOnInit() {
        this.socketService.websocket.subscribe(
            msg => {
                let data = JSON.parse(msg.data);
                switch (data.Type) {
                    case SocketType.CURRENT_CHECKS:
                        this.checks = data.Message;
                        break;
                }
            }
        );
    }

    public startChecks() {
        this.socketService.websocket.next({Type: SocketType.START_CHECKS, Message: this.checks});
    }

    public stopChecks() {
        this.socketService.websocket.next({Type: SocketType.STOP_CHECKS});
    }

    private getCurrentChecks() {
        this.socketService.websocket.next({Type: SocketType.CURRENT_CHECKS});
    }
}
