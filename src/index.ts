import { ECPairFactory } from "ecpair";
import { payments, script, crypto } from "bitcoinjs-lib";
import createHash = require("create-hash");
import tinysecp = require('tiny-secp256k1');

const standardHash = (name: createHash.algorithm, data: Buffer | string) => {
  const h = createHash(name);
  return h.update(data).digest();
}

const hash160 = (data: Buffer | string) => {
  const h1 = standardHash('sha256', data);
  const h2 = standardHash('ripemd160', h1);
  return h2;
};

const hash256= (data: Buffer | string) => {
  const h1 = standardHash('sha256', data);
  const h2 = standardHash('sha256', h1);
  return h2;
}

const s = 'bitcoin is awesome';
console.log('ripemd160 = ' + standardHash('ripemd160', s).toString('hex'));
console.log('  hash160 = ' + hash160(s).toString('hex'));
console.log('   sha256 = ' + standardHash('sha256', s).toString('hex'));
console.log('  hash256 = ' + hash256(s).toString('hex'));

const ECPair = ECPairFactory(tinysecp);
const keyPair = ECPair.makeRandom();
// 打印私钥:
console.log('private key = ' + keyPair.publicKey.toString());
// 以十六进制打印:
console.log('hex = ' + keyPair.publicKey.toString('hex'));
// 补齐32位:
// console.log('hex = ' + keyPair.d.toHex(32));

const wif = 'KwdMAjGmerYanjeui5SHS7JkmpZvVipYvB2LJGU1ZxJwYvP98617'
let ecPair = ECPair.fromWIF(wif);

// 计算公钥:
let pubKey = ecPair.publicKey; // 返回Buffer对象
console.log(pubKey.toString('hex')); // 02或03开头的压缩公钥

const publicKey = '02d0de0aaeaefad02b8bdc8a01a1b8b11c696bd3d66a2c5f10780d95b7df42645c';
ecPair = ECPair.fromPublicKey(Buffer.from(publicKey, 'hex'));
let address = payments.p2pkh({pubkey: ecPair.publicKey}).address; // API 发生了改变
console.log(address);

let
    pubKey1 = '026477115981fe981a6918a6297d9803c4dc04f328f22041bedff886bbc2962e01',
    pubKey2 = '02c96db2302d19b43d4c69368babace7854cc84eb9e061cde51cfa77ca4a22b8b9',
    pubKey3 = '03c6103b3b83e4a24a0e33a4df246ef11772f9992663db0c35759a5e2ebf68d8e9',
    pubkeys = [pubKey1, pubKey2, pubKey3].map(s => Buffer.from(s, 'hex')); // 注意把string转换为Buffer

// 创建2-3 RedeemScript:
const redeemScript = payments.p2wsh({
  redeem: payments.p2ms({m:3, pubkeys})
});
console.log('Redeem script: ' + redeemScript.redeem?.output?.toString('hex'));

// 编码:
// let scriptPubKey =script.scriptHash.output.encode(crypto.hash160(redeemScript));
// address = address.fromOutputScript(scriptPubKey);

console.log('Multisig address: ' + address); // 36NUkt6FWUi3LAWBqWRdDmdTWbt91Yvfu7