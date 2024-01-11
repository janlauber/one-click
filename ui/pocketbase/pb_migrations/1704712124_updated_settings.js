/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("w9y791qv1ymfu3m")

  collection.listRule = "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id ?= blueprint.users.id || @request.auth.id = project.user.id)"
  collection.viewRule = ""

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("w9y791qv1ymfu3m")

  collection.listRule = null
  collection.viewRule = null

  return dao.saveCollection(collection)
})
