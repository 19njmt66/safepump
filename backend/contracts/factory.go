// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// SafePumpFactoryMetaData contains all meta data concerning the SafePumpFactory contract.
var SafePumpFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_uniswapV2Router\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_feeRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"INITIAL_VIRTUAL_ETH_RESERVES\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"INITIAL_VIRTUAL_TOKEN_RESERVES\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_BONDING_CURVE_TOKENS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TARGET_ETH_RAISED\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allTokens\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buy\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"claimCreatorFees\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimPlatformFees\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimablePlatformFees\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createToken\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAmountOutEth\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokensIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAmountOutTokens\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ethIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingCreatorFees\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sell\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFeeRecipient\",\"inputs\":[{\"name\":\"_feeRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUniswapV2Router\",\"inputs\":[{\"name\":\"_router\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tokens\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokensSold\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ethRaised\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"migrated\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uniswapV2Router\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Buy\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Migrated\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"pair\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Sell\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenCreated\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// SafePumpFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use SafePumpFactoryMetaData.ABI instead.
var SafePumpFactoryABI = SafePumpFactoryMetaData.ABI

// SafePumpFactory is an auto generated Go binding around an Ethereum contract.
type SafePumpFactory struct {
	SafePumpFactoryCaller     // Read-only binding to the contract
	SafePumpFactoryTransactor // Write-only binding to the contract
	SafePumpFactoryFilterer   // Log filterer for contract events
}

// SafePumpFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafePumpFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafePumpFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafePumpFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafePumpFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafePumpFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafePumpFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafePumpFactorySession struct {
	Contract     *SafePumpFactory  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafePumpFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafePumpFactoryCallerSession struct {
	Contract *SafePumpFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// SafePumpFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafePumpFactoryTransactorSession struct {
	Contract     *SafePumpFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// SafePumpFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafePumpFactoryRaw struct {
	Contract *SafePumpFactory // Generic contract binding to access the raw methods on
}

// SafePumpFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafePumpFactoryCallerRaw struct {
	Contract *SafePumpFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// SafePumpFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafePumpFactoryTransactorRaw struct {
	Contract *SafePumpFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafePumpFactory creates a new instance of SafePumpFactory, bound to a specific deployed contract.
func NewSafePumpFactory(address common.Address, backend bind.ContractBackend) (*SafePumpFactory, error) {
	contract, err := bindSafePumpFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafePumpFactory{SafePumpFactoryCaller: SafePumpFactoryCaller{contract: contract}, SafePumpFactoryTransactor: SafePumpFactoryTransactor{contract: contract}, SafePumpFactoryFilterer: SafePumpFactoryFilterer{contract: contract}}, nil
}

// NewSafePumpFactoryCaller creates a new read-only instance of SafePumpFactory, bound to a specific deployed contract.
func NewSafePumpFactoryCaller(address common.Address, caller bind.ContractCaller) (*SafePumpFactoryCaller, error) {
	contract, err := bindSafePumpFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafePumpFactoryCaller{contract: contract}, nil
}

// NewSafePumpFactoryTransactor creates a new write-only instance of SafePumpFactory, bound to a specific deployed contract.
func NewSafePumpFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*SafePumpFactoryTransactor, error) {
	contract, err := bindSafePumpFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafePumpFactoryTransactor{contract: contract}, nil
}

// NewSafePumpFactoryFilterer creates a new log filterer instance of SafePumpFactory, bound to a specific deployed contract.
func NewSafePumpFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*SafePumpFactoryFilterer, error) {
	contract, err := bindSafePumpFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafePumpFactoryFilterer{contract: contract}, nil
}

// bindSafePumpFactory binds a generic wrapper to an already deployed contract.
func bindSafePumpFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SafePumpFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafePumpFactory *SafePumpFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafePumpFactory.Contract.SafePumpFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafePumpFactory *SafePumpFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.SafePumpFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafePumpFactory *SafePumpFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.SafePumpFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafePumpFactory *SafePumpFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafePumpFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafePumpFactory *SafePumpFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafePumpFactory *SafePumpFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.contract.Transact(opts, method, params...)
}

// INITIALVIRTUALETHRESERVES is a free data retrieval call binding the contract method 0x89c3d965.
//
// Solidity: function INITIAL_VIRTUAL_ETH_RESERVES() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCaller) INITIALVIRTUALETHRESERVES(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "INITIAL_VIRTUAL_ETH_RESERVES")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITIALVIRTUALETHRESERVES is a free data retrieval call binding the contract method 0x89c3d965.
//
// Solidity: function INITIAL_VIRTUAL_ETH_RESERVES() view returns(uint256)
func (_SafePumpFactory *SafePumpFactorySession) INITIALVIRTUALETHRESERVES() (*big.Int, error) {
	return _SafePumpFactory.Contract.INITIALVIRTUALETHRESERVES(&_SafePumpFactory.CallOpts)
}

// INITIALVIRTUALETHRESERVES is a free data retrieval call binding the contract method 0x89c3d965.
//
// Solidity: function INITIAL_VIRTUAL_ETH_RESERVES() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCallerSession) INITIALVIRTUALETHRESERVES() (*big.Int, error) {
	return _SafePumpFactory.Contract.INITIALVIRTUALETHRESERVES(&_SafePumpFactory.CallOpts)
}

// INITIALVIRTUALTOKENRESERVES is a free data retrieval call binding the contract method 0xe0873234.
//
// Solidity: function INITIAL_VIRTUAL_TOKEN_RESERVES() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCaller) INITIALVIRTUALTOKENRESERVES(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "INITIAL_VIRTUAL_TOKEN_RESERVES")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITIALVIRTUALTOKENRESERVES is a free data retrieval call binding the contract method 0xe0873234.
//
// Solidity: function INITIAL_VIRTUAL_TOKEN_RESERVES() view returns(uint256)
func (_SafePumpFactory *SafePumpFactorySession) INITIALVIRTUALTOKENRESERVES() (*big.Int, error) {
	return _SafePumpFactory.Contract.INITIALVIRTUALTOKENRESERVES(&_SafePumpFactory.CallOpts)
}

// INITIALVIRTUALTOKENRESERVES is a free data retrieval call binding the contract method 0xe0873234.
//
// Solidity: function INITIAL_VIRTUAL_TOKEN_RESERVES() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCallerSession) INITIALVIRTUALTOKENRESERVES() (*big.Int, error) {
	return _SafePumpFactory.Contract.INITIALVIRTUALTOKENRESERVES(&_SafePumpFactory.CallOpts)
}

// MAXBONDINGCURVETOKENS is a free data retrieval call binding the contract method 0x7af92768.
//
// Solidity: function MAX_BONDING_CURVE_TOKENS() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCaller) MAXBONDINGCURVETOKENS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "MAX_BONDING_CURVE_TOKENS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXBONDINGCURVETOKENS is a free data retrieval call binding the contract method 0x7af92768.
//
// Solidity: function MAX_BONDING_CURVE_TOKENS() view returns(uint256)
func (_SafePumpFactory *SafePumpFactorySession) MAXBONDINGCURVETOKENS() (*big.Int, error) {
	return _SafePumpFactory.Contract.MAXBONDINGCURVETOKENS(&_SafePumpFactory.CallOpts)
}

// MAXBONDINGCURVETOKENS is a free data retrieval call binding the contract method 0x7af92768.
//
// Solidity: function MAX_BONDING_CURVE_TOKENS() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCallerSession) MAXBONDINGCURVETOKENS() (*big.Int, error) {
	return _SafePumpFactory.Contract.MAXBONDINGCURVETOKENS(&_SafePumpFactory.CallOpts)
}

// TARGETETHRAISED is a free data retrieval call binding the contract method 0x937136af.
//
// Solidity: function TARGET_ETH_RAISED() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCaller) TARGETETHRAISED(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "TARGET_ETH_RAISED")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TARGETETHRAISED is a free data retrieval call binding the contract method 0x937136af.
//
// Solidity: function TARGET_ETH_RAISED() view returns(uint256)
func (_SafePumpFactory *SafePumpFactorySession) TARGETETHRAISED() (*big.Int, error) {
	return _SafePumpFactory.Contract.TARGETETHRAISED(&_SafePumpFactory.CallOpts)
}

// TARGETETHRAISED is a free data retrieval call binding the contract method 0x937136af.
//
// Solidity: function TARGET_ETH_RAISED() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCallerSession) TARGETETHRAISED() (*big.Int, error) {
	return _SafePumpFactory.Contract.TARGETETHRAISED(&_SafePumpFactory.CallOpts)
}

// AllTokens is a free data retrieval call binding the contract method 0x634282af.
//
// Solidity: function allTokens(uint256 ) view returns(address)
func (_SafePumpFactory *SafePumpFactoryCaller) AllTokens(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "allTokens", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AllTokens is a free data retrieval call binding the contract method 0x634282af.
//
// Solidity: function allTokens(uint256 ) view returns(address)
func (_SafePumpFactory *SafePumpFactorySession) AllTokens(arg0 *big.Int) (common.Address, error) {
	return _SafePumpFactory.Contract.AllTokens(&_SafePumpFactory.CallOpts, arg0)
}

// AllTokens is a free data retrieval call binding the contract method 0x634282af.
//
// Solidity: function allTokens(uint256 ) view returns(address)
func (_SafePumpFactory *SafePumpFactoryCallerSession) AllTokens(arg0 *big.Int) (common.Address, error) {
	return _SafePumpFactory.Contract.AllTokens(&_SafePumpFactory.CallOpts, arg0)
}

// ClaimablePlatformFees is a free data retrieval call binding the contract method 0x0af476cb.
//
// Solidity: function claimablePlatformFees() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCaller) ClaimablePlatformFees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "claimablePlatformFees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ClaimablePlatformFees is a free data retrieval call binding the contract method 0x0af476cb.
//
// Solidity: function claimablePlatformFees() view returns(uint256)
func (_SafePumpFactory *SafePumpFactorySession) ClaimablePlatformFees() (*big.Int, error) {
	return _SafePumpFactory.Contract.ClaimablePlatformFees(&_SafePumpFactory.CallOpts)
}

// ClaimablePlatformFees is a free data retrieval call binding the contract method 0x0af476cb.
//
// Solidity: function claimablePlatformFees() view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCallerSession) ClaimablePlatformFees() (*big.Int, error) {
	return _SafePumpFactory.Contract.ClaimablePlatformFees(&_SafePumpFactory.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_SafePumpFactory *SafePumpFactoryCaller) FeeRecipient(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "feeRecipient")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_SafePumpFactory *SafePumpFactorySession) FeeRecipient() (common.Address, error) {
	return _SafePumpFactory.Contract.FeeRecipient(&_SafePumpFactory.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_SafePumpFactory *SafePumpFactoryCallerSession) FeeRecipient() (common.Address, error) {
	return _SafePumpFactory.Contract.FeeRecipient(&_SafePumpFactory.CallOpts)
}

// GetAllTokens is a free data retrieval call binding the contract method 0x2a5c792a.
//
// Solidity: function getAllTokens() view returns(address[])
func (_SafePumpFactory *SafePumpFactoryCaller) GetAllTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "getAllTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetAllTokens is a free data retrieval call binding the contract method 0x2a5c792a.
//
// Solidity: function getAllTokens() view returns(address[])
func (_SafePumpFactory *SafePumpFactorySession) GetAllTokens() ([]common.Address, error) {
	return _SafePumpFactory.Contract.GetAllTokens(&_SafePumpFactory.CallOpts)
}

// GetAllTokens is a free data retrieval call binding the contract method 0x2a5c792a.
//
// Solidity: function getAllTokens() view returns(address[])
func (_SafePumpFactory *SafePumpFactoryCallerSession) GetAllTokens() ([]common.Address, error) {
	return _SafePumpFactory.Contract.GetAllTokens(&_SafePumpFactory.CallOpts)
}

// GetAmountOutEth is a free data retrieval call binding the contract method 0x69bcc665.
//
// Solidity: function getAmountOutEth(address tokenAddress, uint256 tokensIn) view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCaller) GetAmountOutEth(opts *bind.CallOpts, tokenAddress common.Address, tokensIn *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "getAmountOutEth", tokenAddress, tokensIn)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountOutEth is a free data retrieval call binding the contract method 0x69bcc665.
//
// Solidity: function getAmountOutEth(address tokenAddress, uint256 tokensIn) view returns(uint256)
func (_SafePumpFactory *SafePumpFactorySession) GetAmountOutEth(tokenAddress common.Address, tokensIn *big.Int) (*big.Int, error) {
	return _SafePumpFactory.Contract.GetAmountOutEth(&_SafePumpFactory.CallOpts, tokenAddress, tokensIn)
}

// GetAmountOutEth is a free data retrieval call binding the contract method 0x69bcc665.
//
// Solidity: function getAmountOutEth(address tokenAddress, uint256 tokensIn) view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCallerSession) GetAmountOutEth(tokenAddress common.Address, tokensIn *big.Int) (*big.Int, error) {
	return _SafePumpFactory.Contract.GetAmountOutEth(&_SafePumpFactory.CallOpts, tokenAddress, tokensIn)
}

// GetAmountOutTokens is a free data retrieval call binding the contract method 0x4d6dcae1.
//
// Solidity: function getAmountOutTokens(address tokenAddress, uint256 ethIn) view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCaller) GetAmountOutTokens(opts *bind.CallOpts, tokenAddress common.Address, ethIn *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "getAmountOutTokens", tokenAddress, ethIn)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountOutTokens is a free data retrieval call binding the contract method 0x4d6dcae1.
//
// Solidity: function getAmountOutTokens(address tokenAddress, uint256 ethIn) view returns(uint256)
func (_SafePumpFactory *SafePumpFactorySession) GetAmountOutTokens(tokenAddress common.Address, ethIn *big.Int) (*big.Int, error) {
	return _SafePumpFactory.Contract.GetAmountOutTokens(&_SafePumpFactory.CallOpts, tokenAddress, ethIn)
}

// GetAmountOutTokens is a free data retrieval call binding the contract method 0x4d6dcae1.
//
// Solidity: function getAmountOutTokens(address tokenAddress, uint256 ethIn) view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCallerSession) GetAmountOutTokens(tokenAddress common.Address, ethIn *big.Int) (*big.Int, error) {
	return _SafePumpFactory.Contract.GetAmountOutTokens(&_SafePumpFactory.CallOpts, tokenAddress, ethIn)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SafePumpFactory *SafePumpFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SafePumpFactory *SafePumpFactorySession) Owner() (common.Address, error) {
	return _SafePumpFactory.Contract.Owner(&_SafePumpFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SafePumpFactory *SafePumpFactoryCallerSession) Owner() (common.Address, error) {
	return _SafePumpFactory.Contract.Owner(&_SafePumpFactory.CallOpts)
}

// PendingCreatorFees is a free data retrieval call binding the contract method 0xc598b2f9.
//
// Solidity: function pendingCreatorFees(address ) view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCaller) PendingCreatorFees(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "pendingCreatorFees", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PendingCreatorFees is a free data retrieval call binding the contract method 0xc598b2f9.
//
// Solidity: function pendingCreatorFees(address ) view returns(uint256)
func (_SafePumpFactory *SafePumpFactorySession) PendingCreatorFees(arg0 common.Address) (*big.Int, error) {
	return _SafePumpFactory.Contract.PendingCreatorFees(&_SafePumpFactory.CallOpts, arg0)
}

// PendingCreatorFees is a free data retrieval call binding the contract method 0xc598b2f9.
//
// Solidity: function pendingCreatorFees(address ) view returns(uint256)
func (_SafePumpFactory *SafePumpFactoryCallerSession) PendingCreatorFees(arg0 common.Address) (*big.Int, error) {
	return _SafePumpFactory.Contract.PendingCreatorFees(&_SafePumpFactory.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0xe4860339.
//
// Solidity: function tokens(address ) view returns(address tokenAddress, address creator, uint256 tokensSold, uint256 ethRaised, bool migrated)
func (_SafePumpFactory *SafePumpFactoryCaller) Tokens(opts *bind.CallOpts, arg0 common.Address) (struct {
	TokenAddress common.Address
	Creator      common.Address
	TokensSold   *big.Int
	EthRaised    *big.Int
	Migrated     bool
}, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "tokens", arg0)

	outstruct := new(struct {
		TokenAddress common.Address
		Creator      common.Address
		TokensSold   *big.Int
		EthRaised    *big.Int
		Migrated     bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TokenAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Creator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TokensSold = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.EthRaised = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Migrated = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// Tokens is a free data retrieval call binding the contract method 0xe4860339.
//
// Solidity: function tokens(address ) view returns(address tokenAddress, address creator, uint256 tokensSold, uint256 ethRaised, bool migrated)
func (_SafePumpFactory *SafePumpFactorySession) Tokens(arg0 common.Address) (struct {
	TokenAddress common.Address
	Creator      common.Address
	TokensSold   *big.Int
	EthRaised    *big.Int
	Migrated     bool
}, error) {
	return _SafePumpFactory.Contract.Tokens(&_SafePumpFactory.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0xe4860339.
//
// Solidity: function tokens(address ) view returns(address tokenAddress, address creator, uint256 tokensSold, uint256 ethRaised, bool migrated)
func (_SafePumpFactory *SafePumpFactoryCallerSession) Tokens(arg0 common.Address) (struct {
	TokenAddress common.Address
	Creator      common.Address
	TokensSold   *big.Int
	EthRaised    *big.Int
	Migrated     bool
}, error) {
	return _SafePumpFactory.Contract.Tokens(&_SafePumpFactory.CallOpts, arg0)
}

// UniswapV2Router is a free data retrieval call binding the contract method 0x1694505e.
//
// Solidity: function uniswapV2Router() view returns(address)
func (_SafePumpFactory *SafePumpFactoryCaller) UniswapV2Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SafePumpFactory.contract.Call(opts, &out, "uniswapV2Router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UniswapV2Router is a free data retrieval call binding the contract method 0x1694505e.
//
// Solidity: function uniswapV2Router() view returns(address)
func (_SafePumpFactory *SafePumpFactorySession) UniswapV2Router() (common.Address, error) {
	return _SafePumpFactory.Contract.UniswapV2Router(&_SafePumpFactory.CallOpts)
}

// UniswapV2Router is a free data retrieval call binding the contract method 0x1694505e.
//
// Solidity: function uniswapV2Router() view returns(address)
func (_SafePumpFactory *SafePumpFactoryCallerSession) UniswapV2Router() (common.Address, error) {
	return _SafePumpFactory.Contract.UniswapV2Router(&_SafePumpFactory.CallOpts)
}

// Buy is a paid mutator transaction binding the contract method 0xf088d547.
//
// Solidity: function buy(address tokenAddress) payable returns()
func (_SafePumpFactory *SafePumpFactoryTransactor) Buy(opts *bind.TransactOpts, tokenAddress common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.contract.Transact(opts, "buy", tokenAddress)
}

// Buy is a paid mutator transaction binding the contract method 0xf088d547.
//
// Solidity: function buy(address tokenAddress) payable returns()
func (_SafePumpFactory *SafePumpFactorySession) Buy(tokenAddress common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.Buy(&_SafePumpFactory.TransactOpts, tokenAddress)
}

// Buy is a paid mutator transaction binding the contract method 0xf088d547.
//
// Solidity: function buy(address tokenAddress) payable returns()
func (_SafePumpFactory *SafePumpFactoryTransactorSession) Buy(tokenAddress common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.Buy(&_SafePumpFactory.TransactOpts, tokenAddress)
}

// ClaimCreatorFees is a paid mutator transaction binding the contract method 0x351fee46.
//
// Solidity: function claimCreatorFees() returns()
func (_SafePumpFactory *SafePumpFactoryTransactor) ClaimCreatorFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafePumpFactory.contract.Transact(opts, "claimCreatorFees")
}

// ClaimCreatorFees is a paid mutator transaction binding the contract method 0x351fee46.
//
// Solidity: function claimCreatorFees() returns()
func (_SafePumpFactory *SafePumpFactorySession) ClaimCreatorFees() (*types.Transaction, error) {
	return _SafePumpFactory.Contract.ClaimCreatorFees(&_SafePumpFactory.TransactOpts)
}

// ClaimCreatorFees is a paid mutator transaction binding the contract method 0x351fee46.
//
// Solidity: function claimCreatorFees() returns()
func (_SafePumpFactory *SafePumpFactoryTransactorSession) ClaimCreatorFees() (*types.Transaction, error) {
	return _SafePumpFactory.Contract.ClaimCreatorFees(&_SafePumpFactory.TransactOpts)
}

// ClaimPlatformFees is a paid mutator transaction binding the contract method 0x2bebf6bf.
//
// Solidity: function claimPlatformFees() returns()
func (_SafePumpFactory *SafePumpFactoryTransactor) ClaimPlatformFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafePumpFactory.contract.Transact(opts, "claimPlatformFees")
}

// ClaimPlatformFees is a paid mutator transaction binding the contract method 0x2bebf6bf.
//
// Solidity: function claimPlatformFees() returns()
func (_SafePumpFactory *SafePumpFactorySession) ClaimPlatformFees() (*types.Transaction, error) {
	return _SafePumpFactory.Contract.ClaimPlatformFees(&_SafePumpFactory.TransactOpts)
}

// ClaimPlatformFees is a paid mutator transaction binding the contract method 0x2bebf6bf.
//
// Solidity: function claimPlatformFees() returns()
func (_SafePumpFactory *SafePumpFactoryTransactorSession) ClaimPlatformFees() (*types.Transaction, error) {
	return _SafePumpFactory.Contract.ClaimPlatformFees(&_SafePumpFactory.TransactOpts)
}

// CreateToken is a paid mutator transaction binding the contract method 0x2f2f2d56.
//
// Solidity: function createToken(string name, string symbol) returns(address)
func (_SafePumpFactory *SafePumpFactoryTransactor) CreateToken(opts *bind.TransactOpts, name string, symbol string) (*types.Transaction, error) {
	return _SafePumpFactory.contract.Transact(opts, "createToken", name, symbol)
}

// CreateToken is a paid mutator transaction binding the contract method 0x2f2f2d56.
//
// Solidity: function createToken(string name, string symbol) returns(address)
func (_SafePumpFactory *SafePumpFactorySession) CreateToken(name string, symbol string) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.CreateToken(&_SafePumpFactory.TransactOpts, name, symbol)
}

// CreateToken is a paid mutator transaction binding the contract method 0x2f2f2d56.
//
// Solidity: function createToken(string name, string symbol) returns(address)
func (_SafePumpFactory *SafePumpFactoryTransactorSession) CreateToken(name string, symbol string) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.CreateToken(&_SafePumpFactory.TransactOpts, name, symbol)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SafePumpFactory *SafePumpFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafePumpFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SafePumpFactory *SafePumpFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _SafePumpFactory.Contract.RenounceOwnership(&_SafePumpFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SafePumpFactory *SafePumpFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SafePumpFactory.Contract.RenounceOwnership(&_SafePumpFactory.TransactOpts)
}

// Sell is a paid mutator transaction binding the contract method 0x6c197ff5.
//
// Solidity: function sell(address tokenAddress, uint256 tokenAmount) returns()
func (_SafePumpFactory *SafePumpFactoryTransactor) Sell(opts *bind.TransactOpts, tokenAddress common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _SafePumpFactory.contract.Transact(opts, "sell", tokenAddress, tokenAmount)
}

// Sell is a paid mutator transaction binding the contract method 0x6c197ff5.
//
// Solidity: function sell(address tokenAddress, uint256 tokenAmount) returns()
func (_SafePumpFactory *SafePumpFactorySession) Sell(tokenAddress common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.Sell(&_SafePumpFactory.TransactOpts, tokenAddress, tokenAmount)
}

// Sell is a paid mutator transaction binding the contract method 0x6c197ff5.
//
// Solidity: function sell(address tokenAddress, uint256 tokenAmount) returns()
func (_SafePumpFactory *SafePumpFactoryTransactorSession) Sell(tokenAddress common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.Sell(&_SafePumpFactory.TransactOpts, tokenAddress, tokenAmount)
}

// SetFeeRecipient is a paid mutator transaction binding the contract method 0xe74b981b.
//
// Solidity: function setFeeRecipient(address _feeRecipient) returns()
func (_SafePumpFactory *SafePumpFactoryTransactor) SetFeeRecipient(opts *bind.TransactOpts, _feeRecipient common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.contract.Transact(opts, "setFeeRecipient", _feeRecipient)
}

// SetFeeRecipient is a paid mutator transaction binding the contract method 0xe74b981b.
//
// Solidity: function setFeeRecipient(address _feeRecipient) returns()
func (_SafePumpFactory *SafePumpFactorySession) SetFeeRecipient(_feeRecipient common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.SetFeeRecipient(&_SafePumpFactory.TransactOpts, _feeRecipient)
}

// SetFeeRecipient is a paid mutator transaction binding the contract method 0xe74b981b.
//
// Solidity: function setFeeRecipient(address _feeRecipient) returns()
func (_SafePumpFactory *SafePumpFactoryTransactorSession) SetFeeRecipient(_feeRecipient common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.SetFeeRecipient(&_SafePumpFactory.TransactOpts, _feeRecipient)
}

// SetUniswapV2Router is a paid mutator transaction binding the contract method 0x1419841d.
//
// Solidity: function setUniswapV2Router(address _router) returns()
func (_SafePumpFactory *SafePumpFactoryTransactor) SetUniswapV2Router(opts *bind.TransactOpts, _router common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.contract.Transact(opts, "setUniswapV2Router", _router)
}

// SetUniswapV2Router is a paid mutator transaction binding the contract method 0x1419841d.
//
// Solidity: function setUniswapV2Router(address _router) returns()
func (_SafePumpFactory *SafePumpFactorySession) SetUniswapV2Router(_router common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.SetUniswapV2Router(&_SafePumpFactory.TransactOpts, _router)
}

// SetUniswapV2Router is a paid mutator transaction binding the contract method 0x1419841d.
//
// Solidity: function setUniswapV2Router(address _router) returns()
func (_SafePumpFactory *SafePumpFactoryTransactorSession) SetUniswapV2Router(_router common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.SetUniswapV2Router(&_SafePumpFactory.TransactOpts, _router)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SafePumpFactory *SafePumpFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SafePumpFactory *SafePumpFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.TransferOwnership(&_SafePumpFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SafePumpFactory *SafePumpFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SafePumpFactory.Contract.TransferOwnership(&_SafePumpFactory.TransactOpts, newOwner)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SafePumpFactory *SafePumpFactoryTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafePumpFactory.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SafePumpFactory *SafePumpFactorySession) Receive() (*types.Transaction, error) {
	return _SafePumpFactory.Contract.Receive(&_SafePumpFactory.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SafePumpFactory *SafePumpFactoryTransactorSession) Receive() (*types.Transaction, error) {
	return _SafePumpFactory.Contract.Receive(&_SafePumpFactory.TransactOpts)
}

// SafePumpFactoryBuyIterator is returned from FilterBuy and is used to iterate over the raw logs and unpacked data for Buy events raised by the SafePumpFactory contract.
type SafePumpFactoryBuyIterator struct {
	Event *SafePumpFactoryBuy // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SafePumpFactoryBuyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafePumpFactoryBuy)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SafePumpFactoryBuy)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SafePumpFactoryBuyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafePumpFactoryBuyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafePumpFactoryBuy represents a Buy event raised by the SafePumpFactory contract.
type SafePumpFactoryBuy struct {
	Token       common.Address
	Buyer       common.Address
	TokenAmount *big.Int
	EthAmount   *big.Int
	Fee         *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBuy is a free log retrieval operation binding the contract event 0x00f93dbdb72854b6b6fb35433086556f2635fc83c37080c667496fecfa650fb4.
//
// Solidity: event Buy(address indexed token, address indexed buyer, uint256 tokenAmount, uint256 ethAmount, uint256 fee)
func (_SafePumpFactory *SafePumpFactoryFilterer) FilterBuy(opts *bind.FilterOpts, token []common.Address, buyer []common.Address) (*SafePumpFactoryBuyIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _SafePumpFactory.contract.FilterLogs(opts, "Buy", tokenRule, buyerRule)
	if err != nil {
		return nil, err
	}
	return &SafePumpFactoryBuyIterator{contract: _SafePumpFactory.contract, event: "Buy", logs: logs, sub: sub}, nil
}

// WatchBuy is a free log subscription operation binding the contract event 0x00f93dbdb72854b6b6fb35433086556f2635fc83c37080c667496fecfa650fb4.
//
// Solidity: event Buy(address indexed token, address indexed buyer, uint256 tokenAmount, uint256 ethAmount, uint256 fee)
func (_SafePumpFactory *SafePumpFactoryFilterer) WatchBuy(opts *bind.WatchOpts, sink chan<- *SafePumpFactoryBuy, token []common.Address, buyer []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _SafePumpFactory.contract.WatchLogs(opts, "Buy", tokenRule, buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafePumpFactoryBuy)
				if err := _SafePumpFactory.contract.UnpackLog(event, "Buy", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBuy is a log parse operation binding the contract event 0x00f93dbdb72854b6b6fb35433086556f2635fc83c37080c667496fecfa650fb4.
//
// Solidity: event Buy(address indexed token, address indexed buyer, uint256 tokenAmount, uint256 ethAmount, uint256 fee)
func (_SafePumpFactory *SafePumpFactoryFilterer) ParseBuy(log types.Log) (*SafePumpFactoryBuy, error) {
	event := new(SafePumpFactoryBuy)
	if err := _SafePumpFactory.contract.UnpackLog(event, "Buy", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SafePumpFactoryMigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the SafePumpFactory contract.
type SafePumpFactoryMigratedIterator struct {
	Event *SafePumpFactoryMigrated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SafePumpFactoryMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafePumpFactoryMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SafePumpFactoryMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SafePumpFactoryMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafePumpFactoryMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafePumpFactoryMigrated represents a Migrated event raised by the SafePumpFactory contract.
type SafePumpFactoryMigrated struct {
	Token       common.Address
	Pair        common.Address
	EthAmount   *big.Int
	TokenAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xf2ea3ee6d4d03a11390cbbcd09097d9fe2d7efb1b2825c2b509415d2fb95a7ba.
//
// Solidity: event Migrated(address indexed token, address indexed pair, uint256 ethAmount, uint256 tokenAmount)
func (_SafePumpFactory *SafePumpFactoryFilterer) FilterMigrated(opts *bind.FilterOpts, token []common.Address, pair []common.Address) (*SafePumpFactoryMigratedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}

	logs, sub, err := _SafePumpFactory.contract.FilterLogs(opts, "Migrated", tokenRule, pairRule)
	if err != nil {
		return nil, err
	}
	return &SafePumpFactoryMigratedIterator{contract: _SafePumpFactory.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xf2ea3ee6d4d03a11390cbbcd09097d9fe2d7efb1b2825c2b509415d2fb95a7ba.
//
// Solidity: event Migrated(address indexed token, address indexed pair, uint256 ethAmount, uint256 tokenAmount)
func (_SafePumpFactory *SafePumpFactoryFilterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *SafePumpFactoryMigrated, token []common.Address, pair []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}

	logs, sub, err := _SafePumpFactory.contract.WatchLogs(opts, "Migrated", tokenRule, pairRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafePumpFactoryMigrated)
				if err := _SafePumpFactory.contract.UnpackLog(event, "Migrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMigrated is a log parse operation binding the contract event 0xf2ea3ee6d4d03a11390cbbcd09097d9fe2d7efb1b2825c2b509415d2fb95a7ba.
//
// Solidity: event Migrated(address indexed token, address indexed pair, uint256 ethAmount, uint256 tokenAmount)
func (_SafePumpFactory *SafePumpFactoryFilterer) ParseMigrated(log types.Log) (*SafePumpFactoryMigrated, error) {
	event := new(SafePumpFactoryMigrated)
	if err := _SafePumpFactory.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SafePumpFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SafePumpFactory contract.
type SafePumpFactoryOwnershipTransferredIterator struct {
	Event *SafePumpFactoryOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SafePumpFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafePumpFactoryOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SafePumpFactoryOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SafePumpFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafePumpFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafePumpFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the SafePumpFactory contract.
type SafePumpFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SafePumpFactory *SafePumpFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SafePumpFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SafePumpFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SafePumpFactoryOwnershipTransferredIterator{contract: _SafePumpFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SafePumpFactory *SafePumpFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SafePumpFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SafePumpFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafePumpFactoryOwnershipTransferred)
				if err := _SafePumpFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SafePumpFactory *SafePumpFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*SafePumpFactoryOwnershipTransferred, error) {
	event := new(SafePumpFactoryOwnershipTransferred)
	if err := _SafePumpFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SafePumpFactorySellIterator is returned from FilterSell and is used to iterate over the raw logs and unpacked data for Sell events raised by the SafePumpFactory contract.
type SafePumpFactorySellIterator struct {
	Event *SafePumpFactorySell // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SafePumpFactorySellIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafePumpFactorySell)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SafePumpFactorySell)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SafePumpFactorySellIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafePumpFactorySellIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafePumpFactorySell represents a Sell event raised by the SafePumpFactory contract.
type SafePumpFactorySell struct {
	Token       common.Address
	Seller      common.Address
	TokenAmount *big.Int
	EthAmount   *big.Int
	Fee         *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSell is a free log retrieval operation binding the contract event 0x01fbb57444511e3de5b26ac09ad6bec45c3f9a1e59dd4a0f2b13a240d18476ce.
//
// Solidity: event Sell(address indexed token, address indexed seller, uint256 tokenAmount, uint256 ethAmount, uint256 fee)
func (_SafePumpFactory *SafePumpFactoryFilterer) FilterSell(opts *bind.FilterOpts, token []common.Address, seller []common.Address) (*SafePumpFactorySellIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _SafePumpFactory.contract.FilterLogs(opts, "Sell", tokenRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &SafePumpFactorySellIterator{contract: _SafePumpFactory.contract, event: "Sell", logs: logs, sub: sub}, nil
}

// WatchSell is a free log subscription operation binding the contract event 0x01fbb57444511e3de5b26ac09ad6bec45c3f9a1e59dd4a0f2b13a240d18476ce.
//
// Solidity: event Sell(address indexed token, address indexed seller, uint256 tokenAmount, uint256 ethAmount, uint256 fee)
func (_SafePumpFactory *SafePumpFactoryFilterer) WatchSell(opts *bind.WatchOpts, sink chan<- *SafePumpFactorySell, token []common.Address, seller []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _SafePumpFactory.contract.WatchLogs(opts, "Sell", tokenRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafePumpFactorySell)
				if err := _SafePumpFactory.contract.UnpackLog(event, "Sell", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSell is a log parse operation binding the contract event 0x01fbb57444511e3de5b26ac09ad6bec45c3f9a1e59dd4a0f2b13a240d18476ce.
//
// Solidity: event Sell(address indexed token, address indexed seller, uint256 tokenAmount, uint256 ethAmount, uint256 fee)
func (_SafePumpFactory *SafePumpFactoryFilterer) ParseSell(log types.Log) (*SafePumpFactorySell, error) {
	event := new(SafePumpFactorySell)
	if err := _SafePumpFactory.contract.UnpackLog(event, "Sell", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SafePumpFactoryTokenCreatedIterator is returned from FilterTokenCreated and is used to iterate over the raw logs and unpacked data for TokenCreated events raised by the SafePumpFactory contract.
type SafePumpFactoryTokenCreatedIterator struct {
	Event *SafePumpFactoryTokenCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SafePumpFactoryTokenCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafePumpFactoryTokenCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SafePumpFactoryTokenCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SafePumpFactoryTokenCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafePumpFactoryTokenCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafePumpFactoryTokenCreated represents a TokenCreated event raised by the SafePumpFactory contract.
type SafePumpFactoryTokenCreated struct {
	Token   common.Address
	Creator common.Address
	Name    string
	Symbol  string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTokenCreated is a free log retrieval operation binding the contract event 0xd5d05a8421149c74fd223cfc823befb883babf9bf0b0e4d6bf9c8fdb70e59bb4.
//
// Solidity: event TokenCreated(address indexed token, address indexed creator, string name, string symbol)
func (_SafePumpFactory *SafePumpFactoryFilterer) FilterTokenCreated(opts *bind.FilterOpts, token []common.Address, creator []common.Address) (*SafePumpFactoryTokenCreatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SafePumpFactory.contract.FilterLogs(opts, "TokenCreated", tokenRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &SafePumpFactoryTokenCreatedIterator{contract: _SafePumpFactory.contract, event: "TokenCreated", logs: logs, sub: sub}, nil
}

// WatchTokenCreated is a free log subscription operation binding the contract event 0xd5d05a8421149c74fd223cfc823befb883babf9bf0b0e4d6bf9c8fdb70e59bb4.
//
// Solidity: event TokenCreated(address indexed token, address indexed creator, string name, string symbol)
func (_SafePumpFactory *SafePumpFactoryFilterer) WatchTokenCreated(opts *bind.WatchOpts, sink chan<- *SafePumpFactoryTokenCreated, token []common.Address, creator []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SafePumpFactory.contract.WatchLogs(opts, "TokenCreated", tokenRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafePumpFactoryTokenCreated)
				if err := _SafePumpFactory.contract.UnpackLog(event, "TokenCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTokenCreated is a log parse operation binding the contract event 0xd5d05a8421149c74fd223cfc823befb883babf9bf0b0e4d6bf9c8fdb70e59bb4.
//
// Solidity: event TokenCreated(address indexed token, address indexed creator, string name, string symbol)
func (_SafePumpFactory *SafePumpFactoryFilterer) ParseTokenCreated(log types.Log) (*SafePumpFactoryTokenCreated, error) {
	event := new(SafePumpFactoryTokenCreated)
	if err := _SafePumpFactory.contract.UnpackLog(event, "TokenCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
