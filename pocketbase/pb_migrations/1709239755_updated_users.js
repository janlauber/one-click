/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("_pb_users_auth_")

  collection.listRule = "id = @request.auth.id "
  collection.viewRule = "id = @request.auth.id || @request.auth.id ?~ @collection.blueprints.users.id"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("_pb_users_auth_")

  collection.listRule = "id = @request.auth.id || @request.auth.id ?~ @collection.blueprints.users.id"
  collection.viewRule = ""

  return dao.saveCollection(collection)
})
