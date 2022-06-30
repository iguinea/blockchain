package blockchain

/* LOCK MECHANISM

const cFree = int32(0)
const cLock = int32(1)

func (bc *Blockchain) lockTransactionPool() {
	for !atomic.CompareAndSwapInt32(bc.transactionPoolLock, cFree, cLock) {
		runtime.Gosched()
	}
}
func (bc *Blockchain) unlockTransactionPool() {
	atomic.StoreInt32(bc.transactionPoolLock, cFree)
}

func (bc *Blockchain) lockChain() {
	for !atomic.CompareAndSwapInt32(bc.chainLock, cFree, cLock) {
		runtime.Gosched()
	}
}
func (bc *Blockchain) unlockChain() {
	atomic.StoreInt32(bc.chainLock, cFree)
}
*/
