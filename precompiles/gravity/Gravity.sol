// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.17 .0;

address constant GRAVITY_PRECOMPILE_ADDRESS = 0x1000000000000000000000000000000000000001;

GravityI constant GRAVITY_CONTRACT = GravityI(
    GRAVITY_PRECOMPILE_ADDRESS
);

interface GravityI {
    function sendToEth(
        string memory chainId,
        address destination,
        string memory denom,
        uint256 amount,
        uint256 bridgeFee,
        uint256 chainFee
    ) external returns (bool success);
}
