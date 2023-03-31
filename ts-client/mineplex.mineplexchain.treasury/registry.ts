import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgBurn } from "./types/mineplexchain/treasury/tx";
import { MsgChangeOwner } from "./types/mineplexchain/treasury/tx";
import { MsgMint } from "./types/mineplexchain/treasury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/mineplex.mineplexchain.treasury.MsgBurn", MsgBurn],
    ["/mineplex.mineplexchain.treasury.MsgChangeOwner", MsgChangeOwner],
    ["/mineplex.mineplexchain.treasury.MsgMint", MsgMint],
    
];

export { msgTypes }