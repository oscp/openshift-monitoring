import {Component, OnInit} from '@angular/core';
import {SocketService} from "../socket.service";
import {SocketType} from "../shared/socket.types";

@Component({
    selector: 'app-checks',
    template: `<h4>Checks</h4>
<form #form="ngForm">
    <div class="form-group row">
        <label for="masterapi" class="col-sm-2 col-form-label">Master-API-URLs</label>
        <div class="col-sm-10">
            <input type="text" class="form-control" name="masterapiurl" id="masterapiurl" 
                    [(ngModel)]="checks.MasterApiUrl">
        </div>
    </div>
    <div class="row form-group">
        <div class="form-check form-check-inline">
            <label class="form-check-label">
                <input type="checkbox" name="masterapicheck" id="masterapicheck"
                       [(ngModel)]="checks.MasterApiCheck"> Master-API-Checks
            </label>
        </div>
    </div>

    <div class="form-group row">
        <div class="col-sm-3">
            <button class="btn btn-primary" (click)="startChecks()" *ngIf="!checks.IsRunning">Start Checks
            </button>
            <button class="btn btn-primary" (click)="stopChecks()" *ngIf="checks.IsRunning">Stop Checks
            </button>
        </div>
        <div class="col-sm-2">
            <div class="spinner" *ngIf="checks.IsRunning">
                <div class="rect1"></div>
                <div class="rect2"></div>
                <div class="rect3"></div>
                <div class="rect4"></div>
                <div class="rect5"></div>
            </div>
        </div>
    </div>
</form>
`
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
