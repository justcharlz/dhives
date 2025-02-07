package app

import (
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/cosmos/cosmos-sdk/codec/types"
    "github.com/cosmos/cosmos-sdk/x/auth/tx"
    "github.com/cosmos/cosmos-sdk/client"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
type EncodingConfig struct {
    InterfaceRegistry types.InterfaceRegistry
    Marshaler         codec.Codec
    TxConfig         client.TxConfig
    Amino            *codec.LegacyAmino
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() EncodingConfig {
    amino := codec.NewLegacyAmino()
    interfaceRegistry := types.NewInterfaceRegistry()
    marshaler := codec.NewProtoCodec(interfaceRegistry)
    txConfig := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

    return EncodingConfig{
        InterfaceRegistry: interfaceRegistry,
        Marshaler:        marshaler,
        TxConfig:        txConfig,
        Amino:           amino,
    }
}
