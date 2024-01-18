/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("vs5gr49hpah1g9q")

  collection.listRule = "@request.auth.id != \"\" && (@request.auth.id ?= users.id || @request.auth.id = owner.id)"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("vs5gr49hpah1g9q")

  collection.listRule = "@request.auth.id != \"\" && @request.auth.id ?= users.id || @request.auth.id = owner.id"

  return dao.saveCollection(collection)
})
