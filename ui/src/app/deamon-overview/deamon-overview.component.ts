import {Component, OnInit} from '@angular/core';
import {Subject} from 'rxjs';
import {SocketService} from '../socket.service';
import {SocketType} from '../shared/socket.types';
import {NotificationsService} from 'angular2-notifications';

@Component({
  selector: 'app-deamon-overview',
  template: `
        <h4>Connected Deamons</h4>
        <table class="table table-striped">
            <thead class="thead-inverse">
                <tr>
                <th>Hostname</th>
                <th>Type</th>
                </tr>
            </thead>
            <tbody>
                <tr *ngFor="let d of deamons">
                    <td>{{d.Hostname}}</td>
                    <td>{{d.DeamonType}}</td>
                </tr>
            </tbody>
        </table>
    `
})
export class DeamonOverviewComponent implements OnInit {
  private socket: Subject<any>;
  private deamons: any;

  constructor(private socketService: SocketService, private notificationService: NotificationsService) {
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
            this.deamons = data.Message;
            break;
          case SocketType.NEW_DEAMON.Name:
            this.notificationService.info("Deamon joined", "New deamon joined: " + data.Message);
            this.socket.next({Type: SocketType.ALL_DEAMONS});
            break;
          case SocketType.DEAMON_LEFT.Name:
            this.notificationService.info("Deamon left" , "Deamon left: " + data.Message);
            this.socket.next({Type: SocketType.ALL_DEAMONS});
        }
      }
    );
  }
}
