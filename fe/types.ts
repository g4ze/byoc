

// model Service {
//     id              String   @id @default(cuid())
//     createdAt       DateTime @default(now())
//     name            String
//     arn             String
//     taskFamily      String
//     loadBalancerARN String
//     targetGroupARN  String
//     loadbalancerDNS String
//     desiredCount    Int
//     cluster         String
//     image           String
//     // slug and lbdns is hardcoded into the query
//     // string of the reverse proxy het dns fuctionality
//     slug        String @unique
//     eventListenerARN String
//     user            User     @relation(fields: [userName], references: [userName])
//     userName        String
//     // logging functionality still in development
//     logs            String?
// }


export interface Service{
    id: string;
    createdAt: Date;
    name: string;
    arn: string;
    taskFamily: string;
    loadBalancerARN: string;
    targetGroupARN: string;
    loadbalancerDNS: string;
    desiredCount: number;
    cluster: string;
    image: string;
    slug: string;
    eventListenerARN: string;
    user: string;
    userName: string;
    logs: string;
}
export interface Env {
    [variableName: string]: string;
}
