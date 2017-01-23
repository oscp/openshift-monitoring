import {Component, OnInit} from '@angular/core';
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
export class DeamonsComponent implements OnInit {
  private deamons: any;

  constructor(private socketService: SocketService, private notificationService: NotificationsService) {
    this.getDeamons();
  }

  ngOnInit() {
    this.socketService.websocket.subscribe(
      msg => {
        let data = JSON.parse(msg.data);

        switch (data.WsType) {
          case SocketType.WS_ALL_DEAMONS:
            this.deamons = data.Message;
            break;
          case SocketType.WS_NEW_DEAMON:
            this.notificationService.info("Deamon joined", "New deamon joined: " + data.Message);
            this.getDeamons();
            break;
          case SocketType.WS_DEAMON_LEFT:
            this.notificationService.info("Deamon left", "Deamon left: " + data.Message);
            this.getDeamons();
        }
      }
    );
  }

  getDeamons() {
    this.socketService.websocket.next({WsType: SocketType.WS_ALL_DEAMONS});
  }
}
