import {Component, OnInit, ViewChild} from '@angular/core';
import {NotificationsService} from "angular2-notifications";
import {SocketService} from "../socket.service";
import {SocketType} from "../shared/socket.types";
import {BaseChartDirective} from "ng2-charts";

@Component({
    selector: 'app-results',
    templateUrl: 'results.component.html'
})
export class ResultsComponent implements OnInit {
    // Dognut charts
    public errorLabels: string[] = [];
    public errorData: number[] = [];
    public errorChecks: Array<any> = [];

    public successLabels: string[] = [];
    public successData: number[] = [];
    public successDataBackground: string[] = [];
    public chartType: string = 'doughnut';

    public checkOverviewLabels: string[] = ['Started', 'Finished'];
    public checkOverviewData: number[] = [0, 0];

    // Line Chart
    @ViewChild('lineChart') chart: BaseChartDirective;
    public lineChartType: string = 'line';
    public checkLineData: any = [
        {data: [], label: 'Successful checks'},
        {data: [], label: 'Failed checks'}
    ];
    public checkLineLabels: Array<any> = [];
    public checkLineLegend: boolean = true;
    public checkLineOptions: any = {
        responsive: true
    };
    private LINE_CHART_INTERVAL: number = 5000;
    private lastTime: any;
    private successCount: number = 0;
    private errorCount: number = 0;

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
                    case SocketType.ALL_DEAMONS:
                        this.handleDeamonUpdate(data.Message);
                        break;
                }
            }
        );
    }

    private handleDeamonUpdate(deamons) {
        this.checkOverviewData[0] = 0;
        this.checkOverviewData[1] = 0;

        deamons.forEach(d => {
            this.checkOverviewData[0] += d.StartedChecks;
            this.checkOverviewData[1] += d.FailedChecks + d.SuccessfulChecks;
        });

        // Force UI update
        this.checkOverviewData = this.checkOverviewData.slice();
    }

    private handleResult(msg) {
        // Handle specific by result
        if (msg.IsOk) {
            this.handleSuccessResult(msg);
        } else {
            this.handleErrorResult(msg);
        }

        // Handle Line-Charts
        this.handleLineResult();
    }

    private handleLineResult() {
        let now: any = new Date();
        if (this.lastTime == null || now - this.lastTime > this.LINE_CHART_INTERVAL) {
            // Create a new data point
            this.lastTime = now;
            this.checkLineLabels.push(`${this.lastTime.getHours()}:${this.lastTime.getMinutes()}:${this.lastTime.getSeconds()}`);
            this.checkLineData[0].data.push(this.successCount);
            this.checkLineData[1].data.push(this.errorCount);
            // Cleanup data points if to many
        } else {
            // Add to last data point
            let lastPoint = this.checkLineData[0].data.length - 1;
            this.checkLineData[0].data[lastPoint] += this.successCount;
            this.checkLineData[1].data[lastPoint] += this.errorCount;
        }

        // Cleanup counters
        this.successCount = 0;
        this.errorCount = 0;

        // Update UI because of bug in chartjs:
        this.chart.labels = this.checkLineLabels.slice();
        this.checkLineData = this.checkLineData.slice();
    }

    private handleErrorResult(msg) {
        this.errorCount++;
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

        // Tell the user about it
        msg.Date = new Date();
        this.errorChecks.push(msg);
        this.notificationService.error(`check ${msg.Type} failed.`, msg.Message);
    }

    private handleSuccessResult(msg) {
        this.successCount++;
        let idx = this.successLabels.findIndex(m => m == msg.Type);

        if (idx > -1) {
            this.successData[idx] += 1;
        } else {
            this.successLabels.push(msg.Type);
            this.successData.push(1);
            this.successDataBackground.push(idx % 2 === 0 ? '#36A2EB' : '#FF6384')
        }

        // Enforce refresh
        this.successData = this.successData.slice();
        this.successLabels = this.successLabels.slice();
        this.successDataBackground = this.successDataBackground.slice();
    }
}
