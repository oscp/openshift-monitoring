import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {FormsModule} from '@angular/forms';
import {HttpModule} from '@angular/http';
import {AppComponent} from './app.component';
import {SocketService} from './socket.service';
import {SimpleNotificationsModule} from 'angular2-notifications';
import {JobsComponent} from './job/jobs.component';
import {DeamonsComponent} from './deamon-overview/deamons.component';

@NgModule({
  declarations: [
    AppComponent,
    DeamonsComponent,
    JobsComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    SimpleNotificationsModule
  ],
  providers: [SocketService],
  bootstrap: [AppComponent]
})
export class AppModule {
}
