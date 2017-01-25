import {Component, OnInit, ChangeDetectorRef} from '@angular/core';
import {NotificationsService} from "angular2-notifications";
import {SocketService} from "../socket.service";
import {SocketType} from "../shared/socket.types";

@Component({
    selector: 'app-results',
    template: `<br/>
<h4>Results</h4>
<div class="container-fluid">
<div class="row">
    <div class="col-xs-6">
        <h5>Successful checks</h5>
        <canvas baseChart
                [data]="successData"
                [labels]="successLabels"
                [chartType]="chartType"></canvas>
    </div>
    <div class="col-xs-6">
        <h5>Failed checks</h5>
         <canvas baseChart
                [data]="errorData"
                [labels]="errorLabels"
                [chartType]="chartType"></canvas>
    </div>
</div>
</div>
    `
})
export class ResultsComponent implements OnInit {
    public errorLabels: string[] = [];
    public errorData: number[] = [];

    public successLabels: string[] = [];
    public successData: number[] = [];
    public chartType: string = 'doughnut';

    constructor(private socketService: SocketService, private notificationService: NotificationsService) {
    }

    ngOnInit() {
        this.socketService.websocket.subscribe(
            msg => {
                let data = JSON.parse(msg.data);
                switch (data.Type) {
                    case SocketType.CHECK_RESULT:
                        this.handleResult(data.Message);
                        break;
                }
            }
        );
    }

    private handleResult(msg) {
        if (msg.IsOk) {
            let idx = this.successLabels.findIndex(m => m == msg.Type);

            if (idx > -1) {
                this.successData[idx] += 1;
            } else {
                this.successLabels.push(msg.Type);
                this.successData.push(1);
            }

            // Enforce refresh
            this.successData = this.successData.slice();
            this.successLabels = this.successLabels.slice();
        } else {
            let idx = this.errorLabels.findIndex(m => m == msg.Type);

            if (idx > -1) {
                this.errorData[idx] += 1;
            } else {
                this.errorLabels.push(msg.Type);
                this.errorData.push(1);
            }

            // Enforce refresh
            this.errorData = this.errorData.slice();
            this.errorLabels = this.errorLabels.slice();

            this.notificationService.error(`check ${msg.Type} failed.`, msg.Message);
        }
    }
}
