import { OnGatewayConnection, OnGatewayDisconnect, OnGatewayInit } from '@nestjs/websockets';
import { Prisma } from '@prisma/client';
import { Server, Socket } from 'socket.io';
import { AppService } from 'src/app.service';
export declare class AppGateway implements OnGatewayInit, OnGatewayConnection, OnGatewayDisconnect {
    private appService;
    constructor(appService: AppService);
    server: Server;
    handleSendMessage(client: Socket, payload: Prisma.ChatCreateInput): Promise<void>;
    afterInit(server: any): void;
    handleConnection(client: Socket): void;
    handleDisconnect(client: Socket): void;
}
