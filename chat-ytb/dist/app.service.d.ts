import { PrismaService } from './prisma.service';
import { Chat, Prisma } from '@prisma/client';
export declare class AppService {
    private readonly prisma;
    constructor(prisma: PrismaService);
    createMessage(data: Prisma.ChatCreateInput): Promise<Chat>;
    getMessages(): Promise<Chat[]>;
}
