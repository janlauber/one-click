/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s")

  collection.listRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.viewRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.createRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.updateRule = "@request.auth.id != \"\" && @request.auth.id = user.id"
  collection.deleteRule = "@request.auth.id != \"\" && @request.auth.id = user.id"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s")

  collection.listRule = null
  collection.viewRule = null
  collection.createRule = null
  collection.updateRule = null
  collection.deleteRule = null

  return dao.saveCollection(collection)
})
