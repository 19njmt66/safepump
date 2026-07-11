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

// SafePumpTokenMetaData contains all meta data concerning the SafePumpToken contract.
var SafePumpTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_creator\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CREATOR_ALLOCATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_WALLET_LIMIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SELL_LIMIT_24H\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TOTAL_SUPPLY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VESTING_RELEASE_RATE_PER_DAY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"creator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"factory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockedCreatorAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"incubationComplete\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isExcludedFromMaxWallet\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isExcludedFromSellLimit\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastBuyBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"launchBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrationComplete\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrationTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sellRecords\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amountSold\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastWindowStart\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setIncubationComplete\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMigrationComplete\",\"inputs\":[{\"name\":\"_pair\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_router\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uniswapV2Pair\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"uniswapV2Router\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IncubationCompleted\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MigrationCompleted\",\"inputs\":[{\"name\":\"pair\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"router\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// SafePumpTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use SafePumpTokenMetaData.ABI instead.
var SafePumpTokenABI = SafePumpTokenMetaData.ABI

// SafePumpToken is an auto generated Go binding around an Ethereum contract.
type SafePumpToken struct {
	SafePumpTokenCaller     // Read-only binding to the contract
	SafePumpTokenTransactor // Write-only binding to the contract
	SafePumpTokenFilterer   // Log filterer for contract events
}

// SafePumpTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafePumpTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafePumpTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafePumpTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafePumpTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafePumpTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafePumpTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafePumpTokenSession struct {
	Contract     *SafePumpToken    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafePumpTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafePumpTokenCallerSession struct {
	Contract *SafePumpTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SafePumpTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafePumpTokenTransactorSession struct {
	Contract     *SafePumpTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SafePumpTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafePumpTokenRaw struct {
	Contract *SafePumpToken // Generic contract binding to access the raw methods on
}

// SafePumpTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafePumpTokenCallerRaw struct {
	Contract *SafePumpTokenCaller // Generic read-only contract binding to access the raw methods on
}

// SafePumpTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafePumpTokenTransactorRaw struct {
	Contract *SafePumpTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafePumpToken creates a new instance of SafePumpToken, bound to a specific deployed contract.
func NewSafePumpToken(address common.Address, backend bind.ContractBackend) (*SafePumpToken, error) {
	contract, err := bindSafePumpToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafePumpToken{SafePumpTokenCaller: SafePumpTokenCaller{contract: contract}, SafePumpTokenTransactor: SafePumpTokenTransactor{contract: contract}, SafePumpTokenFilterer: SafePumpTokenFilterer{contract: contract}}, nil
}

// NewSafePumpTokenCaller creates a new read-only instance of SafePumpToken, bound to a specific deployed contract.
func NewSafePumpTokenCaller(address common.Address, caller bind.ContractCaller) (*SafePumpTokenCaller, error) {
	contract, err := bindSafePumpToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafePumpTokenCaller{contract: contract}, nil
}

// NewSafePumpTokenTransactor creates a new write-only instance of SafePumpToken, bound to a specific deployed contract.
func NewSafePumpTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*SafePumpTokenTransactor, error) {
	contract, err := bindSafePumpToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafePumpTokenTransactor{contract: contract}, nil
}

// NewSafePumpTokenFilterer creates a new log filterer instance of SafePumpToken, bound to a specific deployed contract.
func NewSafePumpTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*SafePumpTokenFilterer, error) {
	contract, err := bindSafePumpToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafePumpTokenFilterer{contract: contract}, nil
}

// bindSafePumpToken binds a generic wrapper to an already deployed contract.
func bindSafePumpToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SafePumpTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafePumpToken *SafePumpTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafePumpToken.Contract.SafePumpTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafePumpToken *SafePumpTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafePumpToken.Contract.SafePumpTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafePumpToken *SafePumpTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafePumpToken.Contract.SafePumpTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafePumpToken *SafePumpTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafePumpToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafePumpToken *SafePumpTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafePumpToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafePumpToken *SafePumpTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafePumpToken.Contract.contract.Transact(opts, method, params...)
}

// CREATORALLOCATION is a free data retrieval call binding the contract method 0x2bc73063.
//
// Solidity: function CREATOR_ALLOCATION() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) CREATORALLOCATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "CREATOR_ALLOCATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CREATORALLOCATION is a free data retrieval call binding the contract method 0x2bc73063.
//
// Solidity: function CREATOR_ALLOCATION() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) CREATORALLOCATION() (*big.Int, error) {
	return _SafePumpToken.Contract.CREATORALLOCATION(&_SafePumpToken.CallOpts)
}

// CREATORALLOCATION is a free data retrieval call binding the contract method 0x2bc73063.
//
// Solidity: function CREATOR_ALLOCATION() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) CREATORALLOCATION() (*big.Int, error) {
	return _SafePumpToken.Contract.CREATORALLOCATION(&_SafePumpToken.CallOpts)
}

// MAXWALLETLIMIT is a free data retrieval call binding the contract method 0xc6cfdc5b.
//
// Solidity: function MAX_WALLET_LIMIT() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) MAXWALLETLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "MAX_WALLET_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXWALLETLIMIT is a free data retrieval call binding the contract method 0xc6cfdc5b.
//
// Solidity: function MAX_WALLET_LIMIT() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) MAXWALLETLIMIT() (*big.Int, error) {
	return _SafePumpToken.Contract.MAXWALLETLIMIT(&_SafePumpToken.CallOpts)
}

// MAXWALLETLIMIT is a free data retrieval call binding the contract method 0xc6cfdc5b.
//
// Solidity: function MAX_WALLET_LIMIT() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) MAXWALLETLIMIT() (*big.Int, error) {
	return _SafePumpToken.Contract.MAXWALLETLIMIT(&_SafePumpToken.CallOpts)
}

// SELLLIMIT24H is a free data retrieval call binding the contract method 0xfbb3e58b.
//
// Solidity: function SELL_LIMIT_24H() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) SELLLIMIT24H(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "SELL_LIMIT_24H")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SELLLIMIT24H is a free data retrieval call binding the contract method 0xfbb3e58b.
//
// Solidity: function SELL_LIMIT_24H() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) SELLLIMIT24H() (*big.Int, error) {
	return _SafePumpToken.Contract.SELLLIMIT24H(&_SafePumpToken.CallOpts)
}

// SELLLIMIT24H is a free data retrieval call binding the contract method 0xfbb3e58b.
//
// Solidity: function SELL_LIMIT_24H() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) SELLLIMIT24H() (*big.Int, error) {
	return _SafePumpToken.Contract.SELLLIMIT24H(&_SafePumpToken.CallOpts)
}

// TOTALSUPPLY is a free data retrieval call binding the contract method 0x902d55a5.
//
// Solidity: function TOTAL_SUPPLY() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) TOTALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "TOTAL_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TOTALSUPPLY is a free data retrieval call binding the contract method 0x902d55a5.
//
// Solidity: function TOTAL_SUPPLY() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) TOTALSUPPLY() (*big.Int, error) {
	return _SafePumpToken.Contract.TOTALSUPPLY(&_SafePumpToken.CallOpts)
}

// TOTALSUPPLY is a free data retrieval call binding the contract method 0x902d55a5.
//
// Solidity: function TOTAL_SUPPLY() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) TOTALSUPPLY() (*big.Int, error) {
	return _SafePumpToken.Contract.TOTALSUPPLY(&_SafePumpToken.CallOpts)
}

// VESTINGRELEASERATEPERDAY is a free data retrieval call binding the contract method 0x0fcbac51.
//
// Solidity: function VESTING_RELEASE_RATE_PER_DAY() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) VESTINGRELEASERATEPERDAY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "VESTING_RELEASE_RATE_PER_DAY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VESTINGRELEASERATEPERDAY is a free data retrieval call binding the contract method 0x0fcbac51.
//
// Solidity: function VESTING_RELEASE_RATE_PER_DAY() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) VESTINGRELEASERATEPERDAY() (*big.Int, error) {
	return _SafePumpToken.Contract.VESTINGRELEASERATEPERDAY(&_SafePumpToken.CallOpts)
}

// VESTINGRELEASERATEPERDAY is a free data retrieval call binding the contract method 0x0fcbac51.
//
// Solidity: function VESTING_RELEASE_RATE_PER_DAY() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) VESTINGRELEASERATEPERDAY() (*big.Int, error) {
	return _SafePumpToken.Contract.VESTINGRELEASERATEPERDAY(&_SafePumpToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _SafePumpToken.Contract.Allowance(&_SafePumpToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _SafePumpToken.Contract.Allowance(&_SafePumpToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _SafePumpToken.Contract.BalanceOf(&_SafePumpToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _SafePumpToken.Contract.BalanceOf(&_SafePumpToken.CallOpts, account)
}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_SafePumpToken *SafePumpTokenCaller) Creator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "creator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_SafePumpToken *SafePumpTokenSession) Creator() (common.Address, error) {
	return _SafePumpToken.Contract.Creator(&_SafePumpToken.CallOpts)
}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_SafePumpToken *SafePumpTokenCallerSession) Creator() (common.Address, error) {
	return _SafePumpToken.Contract.Creator(&_SafePumpToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_SafePumpToken *SafePumpTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_SafePumpToken *SafePumpTokenSession) Decimals() (uint8, error) {
	return _SafePumpToken.Contract.Decimals(&_SafePumpToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_SafePumpToken *SafePumpTokenCallerSession) Decimals() (uint8, error) {
	return _SafePumpToken.Contract.Decimals(&_SafePumpToken.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_SafePumpToken *SafePumpTokenCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_SafePumpToken *SafePumpTokenSession) Factory() (common.Address, error) {
	return _SafePumpToken.Contract.Factory(&_SafePumpToken.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_SafePumpToken *SafePumpTokenCallerSession) Factory() (common.Address, error) {
	return _SafePumpToken.Contract.Factory(&_SafePumpToken.CallOpts)
}

// GetLockedCreatorAmount is a free data retrieval call binding the contract method 0x7a5123d7.
//
// Solidity: function getLockedCreatorAmount() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) GetLockedCreatorAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "getLockedCreatorAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockedCreatorAmount is a free data retrieval call binding the contract method 0x7a5123d7.
//
// Solidity: function getLockedCreatorAmount() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) GetLockedCreatorAmount() (*big.Int, error) {
	return _SafePumpToken.Contract.GetLockedCreatorAmount(&_SafePumpToken.CallOpts)
}

// GetLockedCreatorAmount is a free data retrieval call binding the contract method 0x7a5123d7.
//
// Solidity: function getLockedCreatorAmount() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) GetLockedCreatorAmount() (*big.Int, error) {
	return _SafePumpToken.Contract.GetLockedCreatorAmount(&_SafePumpToken.CallOpts)
}

// IncubationComplete is a free data retrieval call binding the contract method 0xdf26ec01.
//
// Solidity: function incubationComplete() view returns(bool)
func (_SafePumpToken *SafePumpTokenCaller) IncubationComplete(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "incubationComplete")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IncubationComplete is a free data retrieval call binding the contract method 0xdf26ec01.
//
// Solidity: function incubationComplete() view returns(bool)
func (_SafePumpToken *SafePumpTokenSession) IncubationComplete() (bool, error) {
	return _SafePumpToken.Contract.IncubationComplete(&_SafePumpToken.CallOpts)
}

// IncubationComplete is a free data retrieval call binding the contract method 0xdf26ec01.
//
// Solidity: function incubationComplete() view returns(bool)
func (_SafePumpToken *SafePumpTokenCallerSession) IncubationComplete() (bool, error) {
	return _SafePumpToken.Contract.IncubationComplete(&_SafePumpToken.CallOpts)
}

// IsExcludedFromMaxWallet is a free data retrieval call binding the contract method 0x6dd3d39f.
//
// Solidity: function isExcludedFromMaxWallet(address ) view returns(bool)
func (_SafePumpToken *SafePumpTokenCaller) IsExcludedFromMaxWallet(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "isExcludedFromMaxWallet", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExcludedFromMaxWallet is a free data retrieval call binding the contract method 0x6dd3d39f.
//
// Solidity: function isExcludedFromMaxWallet(address ) view returns(bool)
func (_SafePumpToken *SafePumpTokenSession) IsExcludedFromMaxWallet(arg0 common.Address) (bool, error) {
	return _SafePumpToken.Contract.IsExcludedFromMaxWallet(&_SafePumpToken.CallOpts, arg0)
}

// IsExcludedFromMaxWallet is a free data retrieval call binding the contract method 0x6dd3d39f.
//
// Solidity: function isExcludedFromMaxWallet(address ) view returns(bool)
func (_SafePumpToken *SafePumpTokenCallerSession) IsExcludedFromMaxWallet(arg0 common.Address) (bool, error) {
	return _SafePumpToken.Contract.IsExcludedFromMaxWallet(&_SafePumpToken.CallOpts, arg0)
}

// IsExcludedFromSellLimit is a free data retrieval call binding the contract method 0x13cda0c8.
//
// Solidity: function isExcludedFromSellLimit(address ) view returns(bool)
func (_SafePumpToken *SafePumpTokenCaller) IsExcludedFromSellLimit(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "isExcludedFromSellLimit", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExcludedFromSellLimit is a free data retrieval call binding the contract method 0x13cda0c8.
//
// Solidity: function isExcludedFromSellLimit(address ) view returns(bool)
func (_SafePumpToken *SafePumpTokenSession) IsExcludedFromSellLimit(arg0 common.Address) (bool, error) {
	return _SafePumpToken.Contract.IsExcludedFromSellLimit(&_SafePumpToken.CallOpts, arg0)
}

// IsExcludedFromSellLimit is a free data retrieval call binding the contract method 0x13cda0c8.
//
// Solidity: function isExcludedFromSellLimit(address ) view returns(bool)
func (_SafePumpToken *SafePumpTokenCallerSession) IsExcludedFromSellLimit(arg0 common.Address) (bool, error) {
	return _SafePumpToken.Contract.IsExcludedFromSellLimit(&_SafePumpToken.CallOpts, arg0)
}

// LastBuyBlock is a free data retrieval call binding the contract method 0x7973bfd2.
//
// Solidity: function lastBuyBlock() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) LastBuyBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "lastBuyBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBuyBlock is a free data retrieval call binding the contract method 0x7973bfd2.
//
// Solidity: function lastBuyBlock() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) LastBuyBlock() (*big.Int, error) {
	return _SafePumpToken.Contract.LastBuyBlock(&_SafePumpToken.CallOpts)
}

// LastBuyBlock is a free data retrieval call binding the contract method 0x7973bfd2.
//
// Solidity: function lastBuyBlock() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) LastBuyBlock() (*big.Int, error) {
	return _SafePumpToken.Contract.LastBuyBlock(&_SafePumpToken.CallOpts)
}

// LaunchBlock is a free data retrieval call binding the contract method 0xd00efb2f.
//
// Solidity: function launchBlock() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) LaunchBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "launchBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LaunchBlock is a free data retrieval call binding the contract method 0xd00efb2f.
//
// Solidity: function launchBlock() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) LaunchBlock() (*big.Int, error) {
	return _SafePumpToken.Contract.LaunchBlock(&_SafePumpToken.CallOpts)
}

// LaunchBlock is a free data retrieval call binding the contract method 0xd00efb2f.
//
// Solidity: function launchBlock() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) LaunchBlock() (*big.Int, error) {
	return _SafePumpToken.Contract.LaunchBlock(&_SafePumpToken.CallOpts)
}

// MigrationComplete is a free data retrieval call binding the contract method 0x2bff884f.
//
// Solidity: function migrationComplete() view returns(bool)
func (_SafePumpToken *SafePumpTokenCaller) MigrationComplete(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "migrationComplete")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MigrationComplete is a free data retrieval call binding the contract method 0x2bff884f.
//
// Solidity: function migrationComplete() view returns(bool)
func (_SafePumpToken *SafePumpTokenSession) MigrationComplete() (bool, error) {
	return _SafePumpToken.Contract.MigrationComplete(&_SafePumpToken.CallOpts)
}

// MigrationComplete is a free data retrieval call binding the contract method 0x2bff884f.
//
// Solidity: function migrationComplete() view returns(bool)
func (_SafePumpToken *SafePumpTokenCallerSession) MigrationComplete() (bool, error) {
	return _SafePumpToken.Contract.MigrationComplete(&_SafePumpToken.CallOpts)
}

// MigrationTime is a free data retrieval call binding the contract method 0xff61a51c.
//
// Solidity: function migrationTime() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) MigrationTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "migrationTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MigrationTime is a free data retrieval call binding the contract method 0xff61a51c.
//
// Solidity: function migrationTime() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) MigrationTime() (*big.Int, error) {
	return _SafePumpToken.Contract.MigrationTime(&_SafePumpToken.CallOpts)
}

// MigrationTime is a free data retrieval call binding the contract method 0xff61a51c.
//
// Solidity: function migrationTime() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) MigrationTime() (*big.Int, error) {
	return _SafePumpToken.Contract.MigrationTime(&_SafePumpToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SafePumpToken *SafePumpTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SafePumpToken *SafePumpTokenSession) Name() (string, error) {
	return _SafePumpToken.Contract.Name(&_SafePumpToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SafePumpToken *SafePumpTokenCallerSession) Name() (string, error) {
	return _SafePumpToken.Contract.Name(&_SafePumpToken.CallOpts)
}

// SellRecords is a free data retrieval call binding the contract method 0x9dd7e695.
//
// Solidity: function sellRecords(address ) view returns(uint128 amountSold, uint128 lastWindowStart)
func (_SafePumpToken *SafePumpTokenCaller) SellRecords(opts *bind.CallOpts, arg0 common.Address) (struct {
	AmountSold      *big.Int
	LastWindowStart *big.Int
}, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "sellRecords", arg0)

	outstruct := new(struct {
		AmountSold      *big.Int
		LastWindowStart *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AmountSold = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LastWindowStart = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SellRecords is a free data retrieval call binding the contract method 0x9dd7e695.
//
// Solidity: function sellRecords(address ) view returns(uint128 amountSold, uint128 lastWindowStart)
func (_SafePumpToken *SafePumpTokenSession) SellRecords(arg0 common.Address) (struct {
	AmountSold      *big.Int
	LastWindowStart *big.Int
}, error) {
	return _SafePumpToken.Contract.SellRecords(&_SafePumpToken.CallOpts, arg0)
}

// SellRecords is a free data retrieval call binding the contract method 0x9dd7e695.
//
// Solidity: function sellRecords(address ) view returns(uint128 amountSold, uint128 lastWindowStart)
func (_SafePumpToken *SafePumpTokenCallerSession) SellRecords(arg0 common.Address) (struct {
	AmountSold      *big.Int
	LastWindowStart *big.Int
}, error) {
	return _SafePumpToken.Contract.SellRecords(&_SafePumpToken.CallOpts, arg0)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SafePumpToken *SafePumpTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SafePumpToken *SafePumpTokenSession) Symbol() (string, error) {
	return _SafePumpToken.Contract.Symbol(&_SafePumpToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SafePumpToken *SafePumpTokenCallerSession) Symbol() (string, error) {
	return _SafePumpToken.Contract.Symbol(&_SafePumpToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SafePumpToken *SafePumpTokenSession) TotalSupply() (*big.Int, error) {
	return _SafePumpToken.Contract.TotalSupply(&_SafePumpToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SafePumpToken *SafePumpTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _SafePumpToken.Contract.TotalSupply(&_SafePumpToken.CallOpts)
}

// UniswapV2Pair is a free data retrieval call binding the contract method 0x49bd5a5e.
//
// Solidity: function uniswapV2Pair() view returns(address)
func (_SafePumpToken *SafePumpTokenCaller) UniswapV2Pair(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "uniswapV2Pair")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UniswapV2Pair is a free data retrieval call binding the contract method 0x49bd5a5e.
//
// Solidity: function uniswapV2Pair() view returns(address)
func (_SafePumpToken *SafePumpTokenSession) UniswapV2Pair() (common.Address, error) {
	return _SafePumpToken.Contract.UniswapV2Pair(&_SafePumpToken.CallOpts)
}

// UniswapV2Pair is a free data retrieval call binding the contract method 0x49bd5a5e.
//
// Solidity: function uniswapV2Pair() view returns(address)
func (_SafePumpToken *SafePumpTokenCallerSession) UniswapV2Pair() (common.Address, error) {
	return _SafePumpToken.Contract.UniswapV2Pair(&_SafePumpToken.CallOpts)
}

// UniswapV2Router is a free data retrieval call binding the contract method 0x1694505e.
//
// Solidity: function uniswapV2Router() view returns(address)
func (_SafePumpToken *SafePumpTokenCaller) UniswapV2Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SafePumpToken.contract.Call(opts, &out, "uniswapV2Router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UniswapV2Router is a free data retrieval call binding the contract method 0x1694505e.
//
// Solidity: function uniswapV2Router() view returns(address)
func (_SafePumpToken *SafePumpTokenSession) UniswapV2Router() (common.Address, error) {
	return _SafePumpToken.Contract.UniswapV2Router(&_SafePumpToken.CallOpts)
}

// UniswapV2Router is a free data retrieval call binding the contract method 0x1694505e.
//
// Solidity: function uniswapV2Router() view returns(address)
func (_SafePumpToken *SafePumpTokenCallerSession) UniswapV2Router() (common.Address, error) {
	return _SafePumpToken.Contract.UniswapV2Router(&_SafePumpToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_SafePumpToken *SafePumpTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _SafePumpToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_SafePumpToken *SafePumpTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _SafePumpToken.Contract.Approve(&_SafePumpToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_SafePumpToken *SafePumpTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _SafePumpToken.Contract.Approve(&_SafePumpToken.TransactOpts, spender, value)
}

// SetIncubationComplete is a paid mutator transaction binding the contract method 0xe6f7ddb1.
//
// Solidity: function setIncubationComplete() returns()
func (_SafePumpToken *SafePumpTokenTransactor) SetIncubationComplete(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafePumpToken.contract.Transact(opts, "setIncubationComplete")
}

// SetIncubationComplete is a paid mutator transaction binding the contract method 0xe6f7ddb1.
//
// Solidity: function setIncubationComplete() returns()
func (_SafePumpToken *SafePumpTokenSession) SetIncubationComplete() (*types.Transaction, error) {
	return _SafePumpToken.Contract.SetIncubationComplete(&_SafePumpToken.TransactOpts)
}

// SetIncubationComplete is a paid mutator transaction binding the contract method 0xe6f7ddb1.
//
// Solidity: function setIncubationComplete() returns()
func (_SafePumpToken *SafePumpTokenTransactorSession) SetIncubationComplete() (*types.Transaction, error) {
	return _SafePumpToken.Contract.SetIncubationComplete(&_SafePumpToken.TransactOpts)
}

// SetMigrationComplete is a paid mutator transaction binding the contract method 0xae49c6d6.
//
// Solidity: function setMigrationComplete(address _pair, address _router) returns()
func (_SafePumpToken *SafePumpTokenTransactor) SetMigrationComplete(opts *bind.TransactOpts, _pair common.Address, _router common.Address) (*types.Transaction, error) {
	return _SafePumpToken.contract.Transact(opts, "setMigrationComplete", _pair, _router)
}

// SetMigrationComplete is a paid mutator transaction binding the contract method 0xae49c6d6.
//
// Solidity: function setMigrationComplete(address _pair, address _router) returns()
func (_SafePumpToken *SafePumpTokenSession) SetMigrationComplete(_pair common.Address, _router common.Address) (*types.Transaction, error) {
	return _SafePumpToken.Contract.SetMigrationComplete(&_SafePumpToken.TransactOpts, _pair, _router)
}

// SetMigrationComplete is a paid mutator transaction binding the contract method 0xae49c6d6.
//
// Solidity: function setMigrationComplete(address _pair, address _router) returns()
func (_SafePumpToken *SafePumpTokenTransactorSession) SetMigrationComplete(_pair common.Address, _router common.Address) (*types.Transaction, error) {
	return _SafePumpToken.Contract.SetMigrationComplete(&_SafePumpToken.TransactOpts, _pair, _router)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_SafePumpToken *SafePumpTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SafePumpToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_SafePumpToken *SafePumpTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SafePumpToken.Contract.Transfer(&_SafePumpToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_SafePumpToken *SafePumpTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SafePumpToken.Contract.Transfer(&_SafePumpToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_SafePumpToken *SafePumpTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SafePumpToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_SafePumpToken *SafePumpTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SafePumpToken.Contract.TransferFrom(&_SafePumpToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_SafePumpToken *SafePumpTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SafePumpToken.Contract.TransferFrom(&_SafePumpToken.TransactOpts, from, to, value)
}

// SafePumpTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the SafePumpToken contract.
type SafePumpTokenApprovalIterator struct {
	Event *SafePumpTokenApproval // Event containing the contract specifics and raw log

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
func (it *SafePumpTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafePumpTokenApproval)
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
		it.Event = new(SafePumpTokenApproval)
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
func (it *SafePumpTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafePumpTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafePumpTokenApproval represents a Approval event raised by the SafePumpToken contract.
type SafePumpTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SafePumpToken *SafePumpTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*SafePumpTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SafePumpToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &SafePumpTokenApprovalIterator{contract: _SafePumpToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SafePumpToken *SafePumpTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *SafePumpTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SafePumpToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafePumpTokenApproval)
				if err := _SafePumpToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SafePumpToken *SafePumpTokenFilterer) ParseApproval(log types.Log) (*SafePumpTokenApproval, error) {
	event := new(SafePumpTokenApproval)
	if err := _SafePumpToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SafePumpTokenIncubationCompletedIterator is returned from FilterIncubationCompleted and is used to iterate over the raw logs and unpacked data for IncubationCompleted events raised by the SafePumpToken contract.
type SafePumpTokenIncubationCompletedIterator struct {
	Event *SafePumpTokenIncubationCompleted // Event containing the contract specifics and raw log

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
func (it *SafePumpTokenIncubationCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafePumpTokenIncubationCompleted)
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
		it.Event = new(SafePumpTokenIncubationCompleted)
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
func (it *SafePumpTokenIncubationCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafePumpTokenIncubationCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafePumpTokenIncubationCompleted represents a IncubationCompleted event raised by the SafePumpToken contract.
type SafePumpTokenIncubationCompleted struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterIncubationCompleted is a free log retrieval operation binding the contract event 0xe6578949dce4d4e246b7a30f2b01dcc2b324089020ec0c1f4746b14fb80155f8.
//
// Solidity: event IncubationCompleted()
func (_SafePumpToken *SafePumpTokenFilterer) FilterIncubationCompleted(opts *bind.FilterOpts) (*SafePumpTokenIncubationCompletedIterator, error) {

	logs, sub, err := _SafePumpToken.contract.FilterLogs(opts, "IncubationCompleted")
	if err != nil {
		return nil, err
	}
	return &SafePumpTokenIncubationCompletedIterator{contract: _SafePumpToken.contract, event: "IncubationCompleted", logs: logs, sub: sub}, nil
}

// WatchIncubationCompleted is a free log subscription operation binding the contract event 0xe6578949dce4d4e246b7a30f2b01dcc2b324089020ec0c1f4746b14fb80155f8.
//
// Solidity: event IncubationCompleted()
func (_SafePumpToken *SafePumpTokenFilterer) WatchIncubationCompleted(opts *bind.WatchOpts, sink chan<- *SafePumpTokenIncubationCompleted) (event.Subscription, error) {

	logs, sub, err := _SafePumpToken.contract.WatchLogs(opts, "IncubationCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafePumpTokenIncubationCompleted)
				if err := _SafePumpToken.contract.UnpackLog(event, "IncubationCompleted", log); err != nil {
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

// ParseIncubationCompleted is a log parse operation binding the contract event 0xe6578949dce4d4e246b7a30f2b01dcc2b324089020ec0c1f4746b14fb80155f8.
//
// Solidity: event IncubationCompleted()
func (_SafePumpToken *SafePumpTokenFilterer) ParseIncubationCompleted(log types.Log) (*SafePumpTokenIncubationCompleted, error) {
	event := new(SafePumpTokenIncubationCompleted)
	if err := _SafePumpToken.contract.UnpackLog(event, "IncubationCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SafePumpTokenMigrationCompletedIterator is returned from FilterMigrationCompleted and is used to iterate over the raw logs and unpacked data for MigrationCompleted events raised by the SafePumpToken contract.
type SafePumpTokenMigrationCompletedIterator struct {
	Event *SafePumpTokenMigrationCompleted // Event containing the contract specifics and raw log

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
func (it *SafePumpTokenMigrationCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafePumpTokenMigrationCompleted)
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
		it.Event = new(SafePumpTokenMigrationCompleted)
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
func (it *SafePumpTokenMigrationCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafePumpTokenMigrationCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafePumpTokenMigrationCompleted represents a MigrationCompleted event raised by the SafePumpToken contract.
type SafePumpTokenMigrationCompleted struct {
	Pair   common.Address
	Router common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMigrationCompleted is a free log retrieval operation binding the contract event 0x8e58716248449261135b3a7b20073d1b4f16376ae21806496b258e9609312693.
//
// Solidity: event MigrationCompleted(address indexed pair, address indexed router)
func (_SafePumpToken *SafePumpTokenFilterer) FilterMigrationCompleted(opts *bind.FilterOpts, pair []common.Address, router []common.Address) (*SafePumpTokenMigrationCompletedIterator, error) {

	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}
	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _SafePumpToken.contract.FilterLogs(opts, "MigrationCompleted", pairRule, routerRule)
	if err != nil {
		return nil, err
	}
	return &SafePumpTokenMigrationCompletedIterator{contract: _SafePumpToken.contract, event: "MigrationCompleted", logs: logs, sub: sub}, nil
}

// WatchMigrationCompleted is a free log subscription operation binding the contract event 0x8e58716248449261135b3a7b20073d1b4f16376ae21806496b258e9609312693.
//
// Solidity: event MigrationCompleted(address indexed pair, address indexed router)
func (_SafePumpToken *SafePumpTokenFilterer) WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *SafePumpTokenMigrationCompleted, pair []common.Address, router []common.Address) (event.Subscription, error) {

	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}
	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _SafePumpToken.contract.WatchLogs(opts, "MigrationCompleted", pairRule, routerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafePumpTokenMigrationCompleted)
				if err := _SafePumpToken.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
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

// ParseMigrationCompleted is a log parse operation binding the contract event 0x8e58716248449261135b3a7b20073d1b4f16376ae21806496b258e9609312693.
//
// Solidity: event MigrationCompleted(address indexed pair, address indexed router)
func (_SafePumpToken *SafePumpTokenFilterer) ParseMigrationCompleted(log types.Log) (*SafePumpTokenMigrationCompleted, error) {
	event := new(SafePumpTokenMigrationCompleted)
	if err := _SafePumpToken.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SafePumpTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the SafePumpToken contract.
type SafePumpTokenTransferIterator struct {
	Event *SafePumpTokenTransfer // Event containing the contract specifics and raw log

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
func (it *SafePumpTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SafePumpTokenTransfer)
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
		it.Event = new(SafePumpTokenTransfer)
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
func (it *SafePumpTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SafePumpTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SafePumpTokenTransfer represents a Transfer event raised by the SafePumpToken contract.
type SafePumpTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SafePumpToken *SafePumpTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SafePumpTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SafePumpToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SafePumpTokenTransferIterator{contract: _SafePumpToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SafePumpToken *SafePumpTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SafePumpTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SafePumpToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SafePumpTokenTransfer)
				if err := _SafePumpToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SafePumpToken *SafePumpTokenFilterer) ParseTransfer(log types.Log) (*SafePumpTokenTransfer, error) {
	event := new(SafePumpTokenTransfer)
	if err := _SafePumpToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
