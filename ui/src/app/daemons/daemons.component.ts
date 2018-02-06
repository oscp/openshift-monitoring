import {Component, OnInit} from '@angular/core';
import {SocketService} from '../socket.service';
import {SocketType} from '../shared/socket.types';
import {NotificationsService} from 'angular2-notifications';

@Component({
  selector: 'app-daemon-overview',
  templateUrl: './daemons.component.html'
})
export class DaemonsComponent implements OnInit {
  private daemons: any;

  constructor(private socketService: SocketService, private notificationService: NotificationsService) {
    this.getDaemons();
  }

  ngOnInit() {
    this.socketService.websocket.subscribe(
      msg => {
        let data = JSON.parse(msg.data);
        switch (data.type) {
          case SocketType.ALL_DAEMONS:
            this.daemons = data.message.sort((a, b) => {
              return a.hostname > b.hostname ? 1 : ((b.hostname > a.hostname) ? -1 : 0);
            });
            break;
          case SocketType.NEW_DAEMON:
            this.notificationService.info('Daemon joined', 'New daemon joined: ' + data.message);
            this.getDaemons();
            break;
          case SocketType.DAEMON_LEFT:
            this.notificationService.info('Daemon left', 'Daemon left: ' + data.message);
            this.getDaemons();
        }
      }
    );
  }

  getDaemons() {
    this.socketService.websocket.next({type: SocketType.ALL_DAEMONS});
  }
}
